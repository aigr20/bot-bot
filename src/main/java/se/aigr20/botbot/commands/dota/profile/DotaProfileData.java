package se.aigr20.botbot.commands.dota.profile;

import java.util.List;

import se.aigr20.botbot.opendota.model.Match;
import se.aigr20.botbot.opendota.model.Peer;
import se.aigr20.botbot.opendota.model.Player;
import se.aigr20.botbot.opendota.model.PlayerHero;
import se.aigr20.botbot.opendota.model.Winrate;

public record DotaProfileData(Player player,
                              List<PlayerHero> heroes,
                              Winrate winrate,
                              List<Match> recentMatches,
                              List<Peer> peers) {
}
