package se.aigr20.botbot;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.concurrent.locks.ReentrantLock;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class BotBotDatasource implements AutoCloseable {
  private static final Logger logger = LoggerFactory.getLogger(BotBotDatasource.class);

  private final Connection connection;
  private final ReentrantLock writeLock = new ReentrantLock();

  public BotBotDatasource(final String dbPath) throws SQLException {
    connection = DriverManager.getConnection("jdbc:sqlite:" + dbPath);
    try (final Statement statement = connection.createStatement()) {
      statement.execute("PRAGMA journal_mode=WAL");
      statement.execute("PRAGMA foreign_keys=ON");
      statement.execute("PRAGMA busy_timeout=5000");
    }

    migrate();
    logger.info("Database opened: {}", dbPath);
  }

  public <T> T write(final SqlTask<T> task) throws SQLException {
    writeLock.lock();
    try {
      return task.run(connection);
    } finally {
      writeLock.unlock();
    }
  }

  public <T> T read(final SqlTask<T> task) throws SQLException {
    return write(task);
  }

  @FunctionalInterface
  public interface SqlTask<T> {
    T run(Connection connection) throws SQLException;
  }

  private void migrate() throws SQLException {
    try (final Statement statement = connection.createStatement()) {
      statement.execute("""
                        CREATE TABLE IF NOT EXISTS accounts (
                          discord_id TEXT NOT NULL PRIMARY KEY,
                          steam_account_id INTEGER NOT NULL
                        )
                        """);
      statement.execute("""
                        CREATE TABLE IF NOT EXISTS heroes (
                          id INTEGER NOT NULL PRIMARY KEY,
                          name TEXT NOT NULL,
                          localized_name TEXT NOT NULL,
                          primary_attribute TEXT NOT NULL,
                          legs INTEGER NOT NULL,
                          attack_type TEXT NOT NULL,
                          roles TEXT NOT NULL
                        )
                        """);
      statement.execute("""
                        CREATE TABLE IF NOT EXISTS hero_aliases (
                          alias TEXT NOT NULL PRIMARY KEY,
                          hero_id INTEGER NOT NULL REFERENCES heroes (id)
                        )
                        """);
    }
  }

  @Override
  public void close() throws SQLException {
    connection.close();
  }
}
