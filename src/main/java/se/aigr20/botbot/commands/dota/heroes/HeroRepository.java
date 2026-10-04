package se.aigr20.botbot.commands.dota.heroes;

import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.Arrays;
import java.util.Optional;

import se.aigr20.botbot.BotBotDatasource;
import se.aigr20.botbot.opendota.model.Hero;

public class HeroRepository {
  private final BotBotDatasource datasource;

  public HeroRepository(final BotBotDatasource datasource) {
    this.datasource = datasource;
  }

  /**
   * Look up the {@link Hero} id matching a hero alias.
   *
   * @param Alias to look up.
   * @return Found hero.
   */
  public Optional<Hero> findHeroByAlias(final String alias) throws SQLException {
    final String sql = """
                       select
                         hero.id, hero.name, hero.localized_name, hero.primary_attribute, hero.legs,
                         hero.roles, hero.attack_type
                       from heroes hero
                       left join hero_aliases hero_alias on hero.id = hero_alias.hero_id
                       where hero_alias.alias = ?
                       """;

    return datasource.read(connection -> {
      try (final PreparedStatement statement = connection.prepareStatement(sql)) {
        statement.setString(1, alias);

        try (final ResultSet result = statement.executeQuery()) {
          if (!result.next()) {
            return Optional.empty();
          }

          final Hero hero = new Hero(result.getInt(1),
                                     result.getString(2),
                                     result.getString(3),
                                     result.getString(4),
                                     result.getString(7),
                                     Arrays.asList(result.getString(6).split(";")),
                                     result.getInt(4));
          return Optional.of(hero);
        }
      }
    });
  }

  public Optional<Hero> findHeroById(final long heroId) throws SQLException {
    final String sql = """
                       select
                         hero.id, hero.name, hero.localized_name, hero.primary_attribute, hero.legs,
                         hero.roles, hero.attack_type
                       from heroes hero
                       where hero.id = ?
                       """;

    return datasource.read(connection -> {
      try (final PreparedStatement statement = connection.prepareStatement(sql)) {
        statement.setLong(1, heroId);

        try (final ResultSet result = statement.executeQuery()) {
          if (!result.next()) {
            return Optional.empty();
          }

          final Hero hero = new Hero(result.getInt(1),
                                     result.getString(2),
                                     result.getString(3),
                                     result.getString(4),
                                     result.getString(7),
                                     Arrays.asList(result.getString(6).split(";")),
                                     result.getInt(5));
          return Optional.of(hero);
        }
      }
    });
  }
}
