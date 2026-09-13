package se.aigr20.botbot.opendota;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.function.Function;
import java.util.stream.Collectors;

import se.aigr20.botbot.opendota.model.Hero;
import se.aigr20.botbot.opendota.model.Match;
import se.aigr20.botbot.opendota.model.TotalField;
import se.aigr20.botbot.opendota.model.Winrate;

import tools.jackson.core.type.TypeReference;
import tools.jackson.databind.ObjectMapper;

public class OpenDotaClient {
  private static final String API_BASE = "https://api.opendota.com/api";

  private final HttpClient http;
  private final ObjectMapper mapper;

  public OpenDotaClient(final ObjectMapper mapper) {
    this.mapper = mapper;
    this.http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(5L)).build();
  }

  public List<Hero> getHeroes() throws OpenDotaException {
    return get("/heroes", new TypeReference<>() {
    });
  }

  /**
   * Get the winrate of a player as a hero.
   *
   * @param accountId Steam account ID of the player.
   * @param heroId    ID of the hero.
   * @return Winrate.
   */
  public Winrate getWinrateAs(final long accountId, final int heroId) throws OpenDotaException {
    return get("/players/%d/wl?hero_id=%d".formatted(accountId, heroId), new TypeReference<>() {
    });
  }

  /**
   * Get the total stats for a player as a hero. The result is converted to a map keyed on the field
   * name.
   *
   * @param accountId Steam account ID of the player.
   * @param heroId    ID of the hero.
   * @return Statistics, by field.
   */
  public Map<String, TotalField> getTotalsAs(final long accountId,
                                             final int heroId) throws OpenDotaException {
    return get("/players/%d/totals?hero_id=%d".formatted(accountId, heroId),
               new TypeReference<List<TotalField>>() {
               }).stream().collect(Collectors.toMap(TotalField::field, Function.identity()));
  }

  public Optional<Match> findLatestMatchAs(final long accountId,
                                           final int heroId) throws OpenDotaException {
    final List<Match> matches = get("/players/%d/matches?hero_id=%d&limit=1&sort=start_time"
            .formatted(accountId, heroId), new TypeReference<>() {
            });

    if (matches.isEmpty()) {
      return Optional.empty();
    }
    return Optional.of(matches.getFirst());
  }

  private <T> T get(final String path, final TypeReference<T> type) throws OpenDotaException {
    final HttpRequest request = HttpRequest
            .newBuilder()
            .uri(URI.create(API_BASE + path))
            .timeout(Duration.ofSeconds(10L))
            .GET()
            .build();

    final HttpResponse<String> response;
    try {
      response = http.send(request, HttpResponse.BodyHandlers.ofString());
    } catch (final IOException e) {
      throw new OpenDotaException("Request failed", e);
    } catch (final InterruptedException e) {
      Thread.currentThread().interrupt();
      throw new OpenDotaException("Request interrupted: " + path, e);
    }

    final int status = response.statusCode();
    if (status != 200) {
      throw new OpenDotaException(status, "OpenDota returned %d for %s".formatted(status, path));
    }

    return mapper.readValue(response.body(), type);
  }
}
