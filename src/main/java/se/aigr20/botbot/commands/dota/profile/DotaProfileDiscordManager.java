
package se.aigr20.botbot.commands.dota.profile;

import java.awt.Color;
import java.math.BigDecimal;
import java.math.RoundingMode;
import java.sql.SQLException;
import java.time.Instant;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.EnumSet;
import java.util.List;
import java.util.Optional;
import java.util.OptionalLong;
import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import se.aigr20.botbot.commands.dota.heroes.HeroRepository;
import se.aigr20.botbot.opendota.OpenDotaException;
import se.aigr20.botbot.opendota.OpenDotaRuntimeException;
import se.aigr20.botbot.opendota.model.Hero;
import se.aigr20.botbot.opendota.model.Match;
import se.aigr20.botbot.opendota.model.Peer;
import se.aigr20.botbot.opendota.model.Player;
import se.aigr20.botbot.opendota.model.PlayerHero;
import se.aigr20.botbot.opendota.model.Profile;

import net.dv8tion.jda.api.EmbedBuilder;
import net.dv8tion.jda.api.components.actionrow.ActionRow;
import net.dv8tion.jda.api.components.buttons.Button;
import net.dv8tion.jda.api.entities.MessageEmbed;
import net.dv8tion.jda.api.entities.emoji.Emoji;
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent;
import net.dv8tion.jda.api.events.interaction.component.ButtonInteractionEvent;
import net.dv8tion.jda.api.utils.messages.MessageEditBuilder;
import net.dv8tion.jda.api.utils.messages.MessageEditData;

public class DotaProfileDiscordManager {
  private static final Logger logger = LoggerFactory.getLogger(DotaProfileDiscordManager.class);
  private static final BigDecimal ONE_HUNDRED = BigDecimal.valueOf(100L);
  private static final DateTimeFormatter DATE_TIME_FORMAT = DateTimeFormatter
          .ofPattern("yyyy-MM-dd HH:mm:ss");

  private final DotaProfileCache profileCache;
  private final HeroRepository heroRepository;

  public DotaProfileDiscordManager(final DotaProfileCache profileCache,
                                   final HeroRepository heroRepository) {
    this.profileCache = profileCache;
    this.heroRepository = heroRepository;
  }

  public void handleButton(final ButtonInteractionEvent event) {
    if (!event.getComponentId().startsWith("dota:profile")) {
      return;
    }

    logger.info("Dota button press: id={}", event.getComponentId());

    event.deferEdit().queue();

    final Tab requestedTab = parseRequestedTab(event)
            .orElseThrow(() -> new IllegalArgumentException("Failed to parse requested tab from button id " +
                                                            event.getComponentId()));
    final long steamId = parseSteamId(event)
            .orElseThrow(() -> new IllegalArgumentException("Failed to parse requested tab from button id " +
                                                            event.getComponentId()));

    final MessageEditData editData = new MessageEditBuilder()
            .setEmbeds(embed(requestedTab, steamId))
            .setComponents(createButtons(requestedTab, steamId))
            .build();

    event.getHook().editOriginal(editData).queue();
  }

  public void handleProfileCommand(final SlashCommandInteractionEvent event, final long steamId) {
    event.deferReply().queue();

    try {
      final MessageEditData editData = new MessageEditBuilder()
              .setEmbeds(embed(Tab.PLAYER, steamId))
              .setComponents(createButtons(Tab.PLAYER, steamId))
              .build();
      event.getHook().editOriginal(editData).queue();
    } catch (final OpenDotaRuntimeException e) {
      final OpenDotaException realException = e.getCause();
      if (realException.isRateLimited()) {
        logger.warn("OpenDota rate limited: account={}", steamId);
        event.getHook().editOriginal("OpenDota is rate limiting us, try again shortly").queue();
        return;
      }

      logger.error("OpenDota request failed: account={}", steamId, realException);
      event.getHook().editOriginal("Failed to retrieve data from OpenDota").queue();
      return;
    }
  }

  private OptionalLong parseSteamId(final ButtonInteractionEvent event) {
    final String[] parts = event.getComponentId().split(":");
    if (parts.length != 4) {
      return OptionalLong.empty();
    }

    try {
      return OptionalLong.of(Long.valueOf(parts[3]));
    } catch (final NumberFormatException e) {
      return OptionalLong.empty();
    }
  }

  private Optional<Tab> parseRequestedTab(final ButtonInteractionEvent event) {
    final String[] parts = event.getComponentId().split(":");
    if (parts.length != 4) {
      return Optional.empty();
    }
    try {
      return Optional.of(Tab.valueOf(parts[2]));
    } catch (final IllegalArgumentException e) {
      return Optional.empty();
    }
  }

  private MessageEmbed[] embed(final Tab tab, final long steamId) {
    final DotaProfileData data = profileCache.get(steamId);

    return switch (tab) {
      case PLAYER -> createPlayerTab(data);
      case PEER -> createPeerTab(data);
      case HISTORY -> createMatchHistoryTab(data);
      case HEROES -> createHeroesTab(data);
    };
  }

  private ActionRow createButtons(final Tab activeTab, final long steamId) {
    final Set<Tab> buttonsToCreate = EnumSet.allOf(Tab.class);
    final List<Button> buttons = new ArrayList<>(buttonsToCreate.size());
    for (final Tab tab : buttonsToCreate) {
      final Button btn = Button
              .secondary("dota:profile:%s:%d".formatted(tab.name(), steamId), tab.getLabel())
              .withEmoji(tab.getEmoji().map(Emoji::fromFormatted).orElse(null))
              .withDisabled(activeTab == tab);
      buttons.add(btn);
    }
    return ActionRow.of(buttons);
  }

  private MessageEmbed[] createPlayerTab(final DotaProfileData data) {
    final Player playerData = data.player();
    final Profile profile = playerData.profile();

    final MessageEmbed embed = getBaseBuilder("Profile of %s".formatted(profile.name()), data)
            .setThumbnail(profile.avatar())
            .addField("Wins", Integer.toString(data.winrate().win()), true)
            .addField("Losses", Integer.toString(data.winrate().lose()), true)
            .build();

    return new MessageEmbed[]{embed};
  }

  private MessageEmbed[] createPeerTab(final DotaProfileData data) {
    final EmbedBuilder embed = getBaseBuilder("Teammates of %s"
            .formatted(data.player().profile().name()), data);

    final List<Peer> peers = data.peers();
    for (int i = 0; i < 5; i++) {
      if (i > peers.size() - 1) {
        break;
      }

      final Peer peer = peers.get(i);

      final BigDecimal wins = BigDecimal.valueOf(peer.gamesWonWith());
      final BigDecimal totalGames = BigDecimal.valueOf(peer.gamesWith());
      final BigDecimal winrate = wins
              .divide(totalGames, 5, RoundingMode.HALF_UP)
              .multiply(ONE_HUNDRED)
              .setScale(2, RoundingMode.HALF_UP);

      embed
              .addField("Name", peer.personaName(), true)
              .addField("Games with", String.valueOf(peer.gamesWith()), true)
              .addField("Winrate with", "%.2f%%".formatted(winrate.doubleValue()), true);
    }

    return new MessageEmbed[]{embed.build()};
  }

  private MessageEmbed[] createMatchHistoryTab(final DotaProfileData data) {
    final List<MessageEmbed> embeds = new ArrayList<>(data.recentMatches().size());

    for (final Match match : data.recentMatches()) {
      final EmbedBuilder embed = new EmbedBuilder();
      final Hero hero;
      try {
        hero = heroRepository
                .findHeroById(match.heroId())
                .orElseThrow(() -> new IllegalArgumentException("Could not find hero with id " +
                                                                match.heroId()));
      } catch (final SQLException e) {
        throw new IllegalStateException("Failed getting hero from database", e);
      }

      embed
              .setTitle("Match %s as %s at %s"
                      .formatted(match.isVictory() ? "won" : "lost",
                                 hero.localizedName(),
                                 Instant
                                         .ofEpochSecond(match.unixStartTime())
                                         .atZone(ZoneId.of("Europe/Stockholm"))
                                         .format(DATE_TIME_FORMAT)))
              .setFooter("Match id: %d".formatted(match.matchId()))
              .setColor(match.isVictory() ? Color.GREEN : Color.RED)
              .addField("Kills", String.valueOf(match.kills()), true)
              .addField("Deaths", String.valueOf(match.deaths()), true)
              .addField("Assists", String.valueOf(match.assists()), true);
      embeds.add(embed.build());
    }

    return embeds.toArray(MessageEmbed[]::new);
  }

  private MessageEmbed[] createHeroesTab(final DotaProfileData data) {
    final EmbedBuilder embed = getBaseBuilder("%s's most played heroes"
            .formatted(data.player().profile().name()), data);

    final List<PlayerHero> heroes = data.heroes();
    for (int i = 0; i < 5; i++) {
      if (i >= heroes.size()) {
        break;
      }

      final PlayerHero hero = heroes.get(i);
      final Hero dotaHero;
      try {
        dotaHero = heroRepository
                .findHeroById(hero.heroId())
                .orElseThrow(() -> new IllegalStateException("Found no hero with id " +
                                                             hero.heroId()));
      } catch (final SQLException e) {
        throw new IllegalStateException("Failed getting hero from database", e);
      }

      final BigDecimal games = BigDecimal.valueOf(hero.games());
      final BigDecimal wins = BigDecimal.valueOf(hero.wins());
      final BigDecimal winrate = wins
              .divide(games, 5, RoundingMode.HALF_UP)
              .multiply(ONE_HUNDRED)
              .setScale(2, RoundingMode.HALF_UP);

      embed
              .addField("Hero", dotaHero.localizedName(), true)
              .addField("Games", String.valueOf(hero.games()), true)
              .addField("Winrate", "%.2f%%".formatted(winrate.doubleValue()), true);
    }

    return new MessageEmbed[]{embed.build()};
  }

  private EmbedBuilder getBaseBuilder(final String title, final DotaProfileData data) {
    final Profile profile = data.player().profile();
    return new EmbedBuilder().setTitle(title).setThumbnail(profile.avatar());
  }

  public static enum Tab {
    PLAYER("Profile"),
    HEROES("Heroes"),
    PEER("Teammates"),
    HISTORY("Recent matches");

    private final String label;
    private final String emoji;

    private Tab(final String label, final String emoji) {
      this.label = label;
      this.emoji = emoji;
    }

    private Tab(final String label) {
      this(label, null);
    }

    public String getLabel() {
      return label;
    }

    public Optional<String> getEmoji() {
      return Optional.ofNullable(emoji);
    }
  }
}
