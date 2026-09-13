package se.aigr20.botbot.opendota.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Winrate(int win, int lose) {
}
