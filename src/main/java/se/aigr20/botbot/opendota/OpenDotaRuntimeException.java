package se.aigr20.botbot.opendota;

public class OpenDotaRuntimeException extends RuntimeException {
  private final OpenDotaException wrappedException;

  public OpenDotaRuntimeException(final OpenDotaException wrapping) {
    super(wrapping);
    this.wrappedException = wrapping;
  }

  @Override
  public OpenDotaException getCause() {
    return wrappedException;
  }
}
