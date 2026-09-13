package se.aigr20.botbot.opendota.model;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Hero(int id,
                   String name,
                   @JsonProperty("localized_name") String localizedName,
                   @JsonProperty("primary_attr") String primaryAttribute,
                   @JsonProperty("attack_type") String attackType,
                   List<String> roles,
                   int legs) {
}
