package se.aigr20.botbot.opendota;

public class OpenDotaException extends Exception {
  private final int status;

  public OpenDotaException(final String message, final Throwable cause) {
    super(message, cause);
    status = -1;
  }

  public OpenDotaException(final int status, final String message) {
    this.status = status;
    super(message);
  }

  /**
   * If the error was not due to the HTTP request receiving an error response, the status will be
   * -1.
   *
   * @return HTTP-status causing the error, or -1 if not applicable.
   */
  public int getStatus() {
    return this.status;
  }

  public boolean isRateLimited() {
    return status == 429;
  }
}
