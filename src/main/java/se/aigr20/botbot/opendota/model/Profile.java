package se.aigr20.botbot.opendota.model;

import java.time.OffsetDateTime;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Profile(@JsonProperty("account_id") long accountId,
                      @JsonProperty("personaname") String name,
                      @JsonProperty("plus") boolean dotaPlus,
                      Integer cheese,
                      @JsonProperty("steamid") String steamId,
                      String avatar,
                      @JsonProperty("avatarmedium") String avatarMedium,
                      @JsonProperty("avatarfull") String avatarFull,
                      @JsonProperty("profileurl") String profileUrl,
                      @JsonProperty("last_login") OffsetDateTime lastLogin,
                      @JsonProperty("loccountrycode") String countryCode,
                      @JsonProperty("is_contributor") boolean contributor,
                      @JsonProperty("is_subscriber") boolean subscriber) {
}
