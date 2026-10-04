package se.aigr20.botbot.opendota.model;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

public record Player(@JsonProperty("rank_tier") Integer rankTier,
                     @JsonProperty("leaderboard_rank") Integer leaderboardRank,
                     @JsonProperty("computed_mrr") Integer mmr,
                     @JsonProperty("computed_mmr_turbo") Integer turboMmr,
                     List<Alias> aliases,
                     Profile profile) {
}
