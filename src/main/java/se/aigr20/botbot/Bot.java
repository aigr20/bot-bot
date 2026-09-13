package se.aigr20.botbot;

import java.io.InputStream;
import java.util.Arrays;
import java.util.HashSet;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import se.aigr20.botbot.commands.dota.AccountRegistry;
import se.aigr20.botbot.commands.dota.DotaCommand;
import se.aigr20.botbot.commands.dota.heroes.HeroAliases;
import se.aigr20.botbot.commands.dota.heroes.HeroSync;
import se.aigr20.botbot.opendota.OpenDotaClient;

import net.dv8tion.jda.api.JDA;
import net.dv8tion.jda.api.JDABuilder;
import net.dv8tion.jda.api.entities.Guild;
import tools.jackson.databind.ObjectMapper;

public class Bot {
  private static final Logger logger = LoggerFactory.getLogger(Bot.class);

  public static void main(final String[] args) throws Exception {
    final ObjectMapper mapper = new ObjectMapper();
    final OpenDotaClient openDotaClient = new OpenDotaClient(mapper);
    final BotBotDatasource datasource = new BotBotDatasource("botbot.db");
    if (args.length > 0 && new HashSet<>(Arrays.asList(args)).contains("sync-heroes")) {
      final HeroSync syncer = new HeroSync(datasource, openDotaClient);
      syncer.syncHeroes();
      try (final InputStream seed = Bot.class.getResourceAsStream("/hero-aliases.txt")) {
        syncer.seedAliases(seed);
      }
      return;
    }

    final String token = System.getenv("DISCORD_TOKEN");
    if (token == null) {
      logger.error("Environment variable DISCORD_TOKEN not set");
      System.exit(1);
    }

    final String guild = System.getenv("GUILD");

    final AccountRegistry accountRegistry = new AccountRegistry(datasource);
    final HeroAliases heroAliases = new HeroAliases(datasource);

    final JDA jda = JDABuilder
            .createDefault(token)
            .addEventListeners(new DotaCommand(accountRegistry, heroAliases, openDotaClient))
            .build();
    jda.awaitReady();

    registerCommands(jda, guild);
    Runtime.getRuntime().addShutdownHook(new Thread(jda::shutdown, "shutdown-jda"));

    logger.info("Bot started as {}", jda.getSelfUser().getName());
  }

  private static void registerCommands(final JDA jda, final String guildId) {
    if (guildId != null) {
      final Guild guild = jda.getGuildById(guildId);
      if (guild == null) {
        logger.error("Configured guild {} not found", guildId);
        return;
      }

      guild
              .updateCommands()
              .addCommands(DotaCommand.specification())
              .queue(_ -> logger.info("Commands registered to guild {}", guildId),
                     e -> logger.error("Failed to register commands", e));
      return;
    }

    jda
            .updateCommands()
            .addCommands(DotaCommand.specification())
            .queue(_ -> logger.info("Commands registered globally"),
                   e -> logger.error("Failed to register commands", e));
  }
}
