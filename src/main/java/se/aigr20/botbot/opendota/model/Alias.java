package se.aigr20.botbot.opendota.model;

import java.time.OffsetDateTime;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Alias(@JsonProperty("personaname") String name,
                    @JsonProperty("name_since") OffsetDateTime nameSince) {
}
