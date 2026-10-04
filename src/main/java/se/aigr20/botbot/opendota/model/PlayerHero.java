package se.aigr20.botbot.opendota.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record PlayerHero(@JsonProperty("hero_id") long heroId,
                         @JsonProperty("last_played") long lastPlayedUnixTimestamp,
                         int games,
                         @JsonProperty("win") int wins,
                         @JsonProperty("with_games") int gamesWith,
                         @JsonProperty("with_win") int gamesWonWith,
                         @JsonProperty("against_games") int gamesAgainst,
                         @JsonProperty("against_win") int gamesWonAgainst) {
}
