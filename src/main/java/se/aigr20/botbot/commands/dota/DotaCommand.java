package se.aigr20.botbot.commands.dota;

import java.sql.SQLException;
import java.util.Objects;
import java.util.Optional;
import java.util.OptionalLong;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import se.aigr20.botbot.commands.dota.heroes.HeroAliases;
import se.aigr20.botbot.opendota.OpenDotaClient;
import se.aigr20.botbot.opendota.OpenDotaException;
import se.aigr20.botbot.opendota.model.Hero;

import net.dv8tion.jda.api.entities.MessageEmbed;
import net.dv8tion.jda.api.entities.User;
import net.dv8tion.jda.api.events.interaction.command.SlashCommandInteractionEvent;
import net.dv8tion.jda.api.hooks.ListenerAdapter;
import net.dv8tion.jda.api.interactions.commands.OptionType;
import net.dv8tion.jda.api.interactions.commands.build.Commands;
import net.dv8tion.jda.api.interactions.commands.build.SlashCommandData;
import net.dv8tion.jda.api.interactions.commands.build.SubcommandData;

public class DotaCommand extends ListenerAdapter {
  public static final String NAME = "dota";

  private static final Logger logger = LoggerFactory.getLogger(DotaCommand.class);
  private static final long STEAM64_BASE = 76561197960265728L;

  private final AccountRegistry accountRegistry;
  private final HeroAliases heroAliases;
  private final OpenDotaClient opendota;

  public DotaCommand(final AccountRegistry accountRegistry,
                     final HeroAliases heroAliases,
                     final OpenDotaClient opendota) {
    this.accountRegistry = accountRegistry;
    this.heroAliases = heroAliases;
    this.opendota = opendota;
  }

  public static SlashCommandData specification() {
    return Commands
            .slash(NAME, "Dota commands")
            .addSubcommands(new SubcommandData("connect",
                                               "Connect a Steam accound to your Discord account")
                                                       .addOption(OptionType.STRING,
                                                                  "steam-id",
                                                                  "Steam account ID (or Steam64 ID)",
                                                                  true),
                            new SubcommandData("stats", "Look up a user's stats for a hero")
                                    .addOption(OptionType.USER, "user", "The user to look up", true)
                                    .addOption(OptionType.STRING,
                                               "hero",
                                               "Hero name, or abbreviation",
                                               true));
  }

  @Override
  public void onSlashCommandInteraction(final SlashCommandInteractionEvent event) {
    if (!NAME.equals(event.getName())) {
      return;
    }

    final String subcommand = event.getSubcommandName();
    logger
            .info("Dota command: subcommand={} user={} ({})",
                  subcommand,
                  event.getUser().getName(),
                  event.getUser().getId());
    Objects.requireNonNull(subcommand);

    switch (subcommand) {
      case "connect" -> handleConnectCommand(event);
      case "stats" -> handleStatsCommand(event);
      default -> event.reply("Unknown subcommand: " + subcommand).setEphemeral(true).queue();
    }
  }

  private void handleConnectCommand(final SlashCommandInteractionEvent event) {
    final String rawId = Objects.requireNonNull(event.getOption("steam-id").getAsString());

    long accountId;
    try {
      accountId = Long.parseLong(rawId);
    } catch (final NumberFormatException e) {
      event.reply("Steam ID must be numeric").setEphemeral(true).queue();
      return;
    }

    if (accountId > STEAM64_BASE) {
      accountId -= STEAM64_BASE;
    }
    if (accountId <= 0) {
      event.reply("That does not seem like a valid Steam ID").setEphemeral(true).queue();
      return;
    }

    try {
      accountRegistry.register(event.getUser().getId(), accountId);
    } catch (final SQLException e) {
      logger.error("Failed to register account: discord_id={}", event.getUser().getId(), e);
      event.reply("").setEphemeral(true).queue();
    }

    logger
            .info("Account registered: discord_id={} steam_account_id={}",
                  event.getUser().getId(),
                  accountId);
    event
            .reply("The Steam account has successfully been connected to your Discord account")
            .queue();
  }

  private void handleStatsCommand(final SlashCommandInteractionEvent event) {
    final User target = Objects.requireNonNull(event.getOption("user").getAsUser());
    final String heroInput = Objects.requireNonNull(event.getOption("hero").getAsString());

    final Optional<Hero> hero;
    try {
      hero = heroAliases.findHeroByAlias(heroInput);
    } catch (final SQLException e) {
      logger.error("Hero lookup failed: hero_alias={}", heroInput, e);
      event.reply("Internal error, please try again later").setEphemeral(true).queue();
      return;
    }

    if (hero.isEmpty()) {
      event
              .reply("Could not find a hero with the alias %s".formatted(heroInput))
              .setEphemeral(true)
              .queue();
      return;
    }

    final OptionalLong steamAccount;
    try {
      steamAccount = accountRegistry.findSteamIdByDiscordId(target.getId());
    } catch (final SQLException e) {
      logger.error("Account lookup failed: discord_id={}", target.getId(), e);
      event.reply("Internal error, please try again later").setEphemeral(true).queue();
      return;
    }

    if (steamAccount.isEmpty()) {
      logger.warn("No user found: account_id={}", target.getId());
      event
              .reply("%s has not connected their Steam account".formatted(target.getAsMention()))
              .queue();
      return;
    }

    event.deferReply().queue();
    try {
      final MessageEmbed embed = new DotaEmbedBuilder(opendota)
              .username(target.getEffectiveName())
              .steamAccount(steamAccount.getAsLong())
              .hero(hero.get())
              .buildStatsEmbed();
      event.getHook().editOriginalEmbeds(embed).queue();
    } catch (final OpenDotaException e) {
      if (e.isRateLimited()) {
        logger.warn("OpenDota rate limited: account={}", steamAccount.getAsLong());
        event.getHook().editOriginal("OpenDota is rate limiting us, try again shortly").queue();
        return;
      }

      logger
              .error("OpenDota request failed: account={} hero={}",
                     steamAccount.getAsLong(),
                     hero.get().id(),
                     e);
      event.getHook().editOriginal("Failed to retrieve data from OpenDota").queue();
    }
  }
}
