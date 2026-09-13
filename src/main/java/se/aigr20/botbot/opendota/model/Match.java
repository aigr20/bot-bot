package se.aigr20.botbot.opendota.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Match(@JsonProperty("match_id") long matchId,
                    @JsonProperty("player_slot") Integer playerSlot,
                    @JsonProperty("radiant_win") Boolean radiantWin,
                    @JsonProperty("duration") int durationSeconds,
                    @JsonProperty("game_mode") int gameMode,
                    @JsonProperty("lobby_type") int lobbyType,
                    @JsonProperty("hero_id") int heroId,
                    @JsonProperty("start_time") long unixStartTime,
                    Integer version,
                    int kills,
                    int deaths,
                    int assists,
                    Integer skill,
                    @JsonProperty("average_rank") Integer averageRank,
                    @JsonProperty("leaver_status") int leaverStatus,
                    @JsonProperty("party_size") Integer partySize,
                    @JsonProperty("hero_variant") Integer heroVariant) {
}
