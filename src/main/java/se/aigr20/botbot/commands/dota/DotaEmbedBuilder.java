package se.aigr20.botbot.commands.dota;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.List;
import java.util.Map;
import java.util.Optional;

import se.aigr20.botbot.opendota.OpenDotaClient;
import se.aigr20.botbot.opendota.OpenDotaException;
import se.aigr20.botbot.opendota.model.Hero;
import se.aigr20.botbot.opendota.model.Match;
import se.aigr20.botbot.opendota.model.TotalField;
import se.aigr20.botbot.opendota.model.Winrate;

import net.dv8tion.jda.api.EmbedBuilder;
import net.dv8tion.jda.api.entities.MessageEmbed;

public class DotaEmbedBuilder {
  private static final BigDecimal ONE_HUNDRED = BigDecimal.valueOf(100);

  private static final List<AverageField> AVERAGE_FIELDS = List
          .of(new AverageField("kills", "Average Kills"),
              new AverageField("deaths", "Average Deaths"),
              new AverageField("assists", "Average Assists"),
              new AverageField("last_hits", "Average Last Hits"),
              new AverageField("denies", "Average Denies"),
              new AverageField("hero_damage", "Average Hero Damage"),
              new AverageField("hero_healing", "Average Hero Healing"),
              new AverageField("purchase_ward_observer", "Avg Observer Wards Bought"),
              new AverageField("purchase_ward_sentry", "Avg Sentry Wards Bought"),
              new AverageField("stuns", "Avg Stun Time Inflicted (s)"));

  private final OpenDotaClient client;
  private String username;
  private Hero hero;
  private long steamAccount;

  public DotaEmbedBuilder(final OpenDotaClient client) {
    this.client = client;
  }

  public DotaEmbedBuilder username(final String username) {
    this.username = username;
    return this;
  }

  public DotaEmbedBuilder hero(final Hero hero) {
    this.hero = hero;
    return this;
  }

  public DotaEmbedBuilder steamAccount(final long steamAccount) {
    this.steamAccount = steamAccount;
    return this;
  }

  public MessageEmbed buildStatsEmbed() throws OpenDotaException {
    final Winrate winrate = client.getWinrateAs(steamAccount, hero.id());
    final Map<String, TotalField> totals = client.getTotalsAs(steamAccount, hero.id());
    final Optional<Match> latestMatch = client.findLatestMatchAs(steamAccount, hero.id());

    final EmbedBuilder embed = new EmbedBuilder()
            .setTitle("Statistics for %s as %s".formatted(username, hero.localizedName()));

    final BigDecimal wins = BigDecimal.valueOf(winrate.win());
    final BigDecimal losses = BigDecimal.valueOf(winrate.lose());
    final BigDecimal games = wins.add(losses);
    final BigDecimal winsShare = wins.divide(games, 5, RoundingMode.HALF_UP);
    final BigDecimal winratePercent = winsShare
            .multiply(ONE_HUNDRED)
            .setScale(2, RoundingMode.HALF_UP);

    embed
            .addField("Wins", String.valueOf(winrate.win()), true)
            .addField("Losses", String.valueOf(winrate.lose()), true)
            .addField("Winrate", "%.2f%".formatted(winratePercent.doubleValue()), true);
    for (final AverageField field : AVERAGE_FIELDS) {
      final TotalField f = totals.get(field.field());
      if (f != null && f.n() > 0) {
        embed.addField("%s".formatted(field.title()), "%.2f".formatted(f.sum() / f.n()), true);
      }
    }

    final String latestValue = latestMatch
            .map(m -> "<t:%d:R>".formatted(m.unixStartTime()))
            .orElse("No games played");
    embed.addField("Latest Game", latestValue, false);

    return embed.build();
  }

  private record AverageField(String field, String title) {
  }
}
