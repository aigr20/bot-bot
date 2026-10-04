package se.aigr20.botbot.commands.dota.profile;

import java.time.Clock;
import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import se.aigr20.botbot.opendota.OpenDotaClient;
import se.aigr20.botbot.opendota.OpenDotaException;
import se.aigr20.botbot.opendota.OpenDotaRuntimeException;
import se.aigr20.botbot.opendota.model.Match;
import se.aigr20.botbot.opendota.model.Peer;
import se.aigr20.botbot.opendota.model.Player;
import se.aigr20.botbot.opendota.model.PlayerHero;
import se.aigr20.botbot.opendota.model.Winrate;

public class DotaProfileCache {
  static final Duration CACHE_TTL = Duration.ofHours(1L);
  static final Duration CLEANUP_INTERVAL = Duration.ofMinutes(10L);

  private final Map<Long, CacheEntry> cache;
  private final OpenDotaClient openDotaClient;
  private final ScheduledExecutorService scheduler;
  private final Clock clock;

  public DotaProfileCache(final OpenDotaClient openDotaClient,
                          final ScheduledExecutorService cleanUpScheduler,
                          final Clock clock) {
    this.cache = new ConcurrentHashMap<>();
    this.scheduler = cleanUpScheduler;
    this.openDotaClient = openDotaClient;
    this.clock = clock;

    this.scheduler
            .scheduleWithFixedDelay(this::removeExpiredEntries,
                                    CLEANUP_INTERVAL.toMinutes(),
                                    CLEANUP_INTERVAL.toMinutes(),
                                    TimeUnit.MINUTES);
  }

  public DotaProfileData get(final long steamAccount) throws OpenDotaRuntimeException {
    final Instant now = clock.instant();

    final CacheEntry entry = cache.compute(steamAccount, (account, cached) -> {
      if (cached == null || cached.expiresAt().isBefore(now)) {
        final DotaProfileData data;
        try {
          data = loadFromApi(account);
        } catch (final OpenDotaException e) {
          throw new OpenDotaRuntimeException(e);
        }
        return new CacheEntry(now.plusSeconds(CACHE_TTL.toSeconds()), data);
      }
      return cached;
    });

    return entry.data();
  }

  public void invalidate(final long steamAccount) {
    cache.remove(steamAccount);
  }

  private void removeExpiredEntries() {
    cache.entrySet().removeIf(entry -> entry.getValue().expiresAt().isBefore(clock.instant()));
  }

  private DotaProfileData loadFromApi(final long steamAccount) throws OpenDotaException {
    final Winrate winrate = openDotaClient.getWinrate(steamAccount);
    final List<Match> recentMatches = openDotaClient.getMatches(steamAccount);
    final List<Peer> peers = openDotaClient.getPeers(steamAccount);
    final Player player = openDotaClient.getPlayer(steamAccount);
    final List<PlayerHero> playerHeroes = openDotaClient.getPlayerHeroes(steamAccount);

    return new DotaProfileData(player, playerHeroes, winrate, recentMatches, peers);
  }

  private record CacheEntry(Instant expiresAt, DotaProfileData data) {
  }
}
