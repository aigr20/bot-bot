package se.aigr20.botbot.commands.dota.profile;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.time.Clock;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;
import org.mockito.Captor;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import se.aigr20.botbot.opendota.OpenDotaClient;
import se.aigr20.botbot.opendota.OpenDotaException;

class DotaProfileCacheTest {
  DotaProfileCache sut;

  @Mock Clock clock;
  @Mock OpenDotaClient openDotaClient;
  @Mock ScheduledExecutorService scheduledExecutorService;

  @Captor ArgumentCaptor<Runnable> cleanUpMethodCaptor;

  AutoCloseable openedMocks;

  final Instant now = LocalDateTime
          .of(2026, 10, 4, 10, 0)
          .atZone(ZoneId.of("Europe/Stockholm"))
          .toInstant();

  @BeforeEach
  void setUp() {
    openedMocks = MockitoAnnotations.openMocks(this);

    sut = new DotaProfileCache(openDotaClient, scheduledExecutorService, clock);

    verify(scheduledExecutorService)
            .scheduleWithFixedDelay(cleanUpMethodCaptor.capture(),
                                    anyLong(),
                                    anyLong(),
                                    any(TimeUnit.class));
  }

  @AfterEach
  void tearDown() throws Exception {
    openedMocks.close();
  }

  @Test
  void when_CachedEntry_then_OnlyGetFromApiOnce() throws OpenDotaException {
    when(clock.instant()).thenReturn(now);

    sut.get(123L);
    sut.get(123L);

    verify(openDotaClient, times(1)).getPlayer(123L);
  }

  @Test
  void when_CachedEntryExpired_then_GetFromApi() throws OpenDotaException {
    when(clock.instant())
            .thenReturn(now,
                        now.plusSeconds(DotaProfileCache.CACHE_TTL.plusMinutes(1L).toSeconds()));

    sut.get(123L);
    sut.get(123L);

    verify(openDotaClient, times(2)).getPlayer(123L);
  }

  @Test
  void when_CachedEntryRemovedByCleanUp_then_GetFromApi() throws OpenDotaException {
    when(clock.instant())
            .thenReturn(now,
                        now.plusSeconds(DotaProfileCache.CACHE_TTL.plusMinutes(1L).toSeconds()),
                        now);

    sut.get(123L);
    cleanUpMethodCaptor.getValue().run();
    sut.get(123L);

    verify(openDotaClient, times(2)).getPlayer(123L);
  }
}
