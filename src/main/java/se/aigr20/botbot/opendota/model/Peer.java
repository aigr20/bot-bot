package se.aigr20.botbot.opendota.model;

import java.time.OffsetDateTime;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Peer(@JsonProperty("account_id") long accountId,
                   @JsonProperty("last_played") long lastPlayed,
                   int wins,
                   int games,
                   @JsonProperty("with_win") int gamesWonWith,
                   @JsonProperty("with_games") int gamesWith,
                   @JsonProperty("against_games") int gamesAgainst,
                   @JsonProperty("against_win") int gamesWonAgainst,
                   @JsonProperty("with_gpm_sum") int gpmSumWith,
                   @JsonProperty("with_xpm_sum") int xpmSumWith,
                   @JsonProperty("personaname") String personaName,
                   String name,
                   @JsonProperty("is_contributor") boolean contributor,
                   @JsonProperty("is_subscriber") boolean subscriber,
                   @JsonProperty("last_login") OffsetDateTime lastLogin,
                   String avatar,
                   @JsonProperty("avatarfull") String fullAvatar) {
}
