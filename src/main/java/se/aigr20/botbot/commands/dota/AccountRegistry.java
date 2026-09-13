package se.aigr20.botbot.commands.dota;

import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.Objects;
import java.util.OptionalLong;

import se.aigr20.botbot.BotBotDatasource;

/**
 * Manages the connection between a Discord account and Steam account.
 */
public class AccountRegistry {
  private final BotBotDatasource datasource;

  public AccountRegistry(final BotBotDatasource datasource) {
    this.datasource = Objects.requireNonNull(datasource);
  }

  /**
   * Connect a Discord account to a Steam account. If the Discord account is already connected to a
   * Steam account, it is replaced.
   *
   * @param discordId The Discord account ID.
   * @param steamId   The Steam account ID.
   */
  public void register(final String discordId, final long steamId) throws SQLException {
    final String sql = """
                       insert into accounts (discord_id, steam_account_id)
                       values (?, ?)
                       on conflict (discord_id)
                       do update set steam_account_id = excluded.steam_account_id
                       """;
    datasource.write(connection -> {
      try (final PreparedStatement stmt = connection.prepareStatement(sql)) {
        stmt.setString(1, discordId);
        stmt.setLong(2, steamId);
        stmt.executeUpdate();
        return null;
      }
    });
  }

  /**
   * Find the Steam account ID connected to the Discord account, if one exists.
   *
   * @param discordId Discord account ID to search for.
   * @return The found Steam account ID.
   */
  public OptionalLong findSteamIdByDiscordId(final String discordId) throws SQLException {
    final String sql = "select steam_account_id from accounts where discord_id = ?";
    return datasource.read(connection -> {
    try (final PreparedStatement stmt = connection.prepareStatement(sql)) {
      stmt.setString(1, discordId);
      try (final ResultSet result = stmt.executeQuery()) {
        return result.next() ? OptionalLong.of(result.getLong(1)) : OptionalLong.empty();
      }
    }
    });
  }
}
