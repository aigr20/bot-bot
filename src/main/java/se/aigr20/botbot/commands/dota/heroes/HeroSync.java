package se.aigr20.botbot.commands.dota.heroes;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.util.List;
import java.util.Locale;
import java.util.stream.Collectors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import se.aigr20.botbot.Bot;
import se.aigr20.botbot.BotBotDatasource;
import se.aigr20.botbot.opendota.OpenDotaClient;
import se.aigr20.botbot.opendota.OpenDotaException;
import se.aigr20.botbot.opendota.model.Hero;

public class HeroSync {
  private static final Logger logger = LoggerFactory.getLogger(HeroSync.class);

  private final BotBotDatasource datasource;
  private final OpenDotaClient openDotaClient;

  public HeroSync(final BotBotDatasource datasource, final OpenDotaClient openDotaClient) {
    this.datasource = datasource;
    this.openDotaClient = openDotaClient;
  }

  public void syncHeroes() throws OpenDotaException, SQLException {
    final List<Hero> heroes = openDotaClient.getHeroes();

    final String sql = """
                       insert into heroes (id, name, localized_name, primary_attribute, legs, attack_type, roles)
                       values (?, ?, ?, ?, ?, ?, ?)
                       on conflict (id)
                       do update set
                         name = excluded.name,
                         localized_name = excluded.localized_name,
                         primary_attribute = excluded.primary_attribute,
                         attack_type = excluded.attack_type,
                         legs = excluded.legs,
                         roles = excluded.roles
                       """;
    datasource.write(connection -> {
      connection.setAutoCommit(false);
      try (final PreparedStatement statement = connection.prepareStatement(sql)) {
        for (final Hero hero : heroes) {
          int p = 1;
          statement.setInt(p++, hero.id());
          statement.setString(p++, hero.name());
          statement.setString(p++, hero.localizedName());
          statement.setString(p++, hero.primaryAttribute());
          statement.setInt(p++, hero.legs());
          statement.setString(p++, hero.attackType());
          statement.setString(p++, hero.roles().stream().collect(Collectors.joining(";")));
          statement.addBatch();
        }

        statement.executeBatch();
        connection.commit();
      } catch (final SQLException e) {
        connection.rollback();
        throw e;
      } finally {
        connection.setAutoCommit(true);
      }
      return null;
    });

    logger.info("Synced {} heroes", heroes.size());
  }

  public void seedAliases(final InputStream seedFile) throws IOException, SQLException {
    final String sql = """
                       insert into hero_aliases (alias, hero_id)
                       select ?, id from heroes where id = ?
                       on conflict (alias)
                       do update set hero_id = excluded.hero_id
                       """;

    datasource.write(connection -> {
      try (final BufferedReader reader = new BufferedReader(new InputStreamReader(seedFile));
           final PreparedStatement statement = connection.prepareStatement(sql)) {
       String line;
        while ((line = reader.readLine()) != null) {
          line = line.trim();
          if (line.isEmpty() || line.startsWith("#")) {
            continue;
          }

          final String[] parts = line.split(":", 2);
          final String hero = parts[0].trim();
          final String[] aliases = parts[1].split(",");
          for (final String alias : aliases) {
            statement.setString(1, alias.trim().toLowerCase(Locale.ROOT));
            statement.setInt(2, Integer.parseInt(hero));
            statement.addBatch();
          }
        }

        statement.executeBatch();
      } catch (final IOException e) {
        throw new IllegalStateException(e);
      }
      return null;
    });

    logger.info("Seeded aliases");
  }
}
