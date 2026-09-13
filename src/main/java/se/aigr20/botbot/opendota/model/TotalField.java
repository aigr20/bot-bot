package se.aigr20.botbot.opendota.model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TotalField(String field, int n, double sum) {
}
