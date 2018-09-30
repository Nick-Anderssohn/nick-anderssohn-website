library error;

enum Level { critical, error, opsError, warning, info, debug }

class LeveledException implements Exception, UserFacingMessage {
  Level _level;

  Level get level => _level;
  final String message;
  String userFacingMessage;
  final List<LeveledException> nestedExceptions;
  StackTrace attachedStackTrace; // Will be null if this isn't a nested exception

  LeveledException(this.message, this._level, {this.nestedExceptions = const []});

  LeveledException.Critical(this.message, {this.nestedExceptions = const []}) {
    _level = Level.critical;
  }

  LeveledException.Error(this.message, {this.nestedExceptions = const []}) {
    _level = Level.error;
  }

  LeveledException.OpsError(this.message, {this.nestedExceptions = const []}) {
    _level = Level.opsError;
  }

  LeveledException.Warning(this.message, {this.nestedExceptions = const []}) {
    _level = Level.warning;
  }

  LeveledException.Info(this.message, {this.nestedExceptions = const []}) {
    _level = Level.info;
  }

  LeveledException.Debug(this.message, {this.nestedExceptions = const []}) {
    _level = Level.debug;
  }

  @override
  String toString() => message == null ? '' : message;

  @override
  String getUserFacingMessage() => userFacingMessage == null ? '' : userFacingMessage;

  @override
  bool hasUserFacingMessage() => userFacingMessage != null;
}

abstract class UserFacingMessage {
  String getUserFacingMessage();
  bool hasUserFacingMessage();
}
