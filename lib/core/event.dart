import 'dart:async';

import 'package:fl_clash/enum/enum.dart';
import 'package:fl_clash/models/models.dart';
import 'package:flutter/foundation.dart';

abstract mixin class CoreEventListener {
  void onLog(Log log) {}

  void onDelay(Delay delay) {}

  void onRequest(TrackerInfo connection) {}

  void onLoaded(String providerName) {}

  void onCrash(String message) {}
}

class CoreEventManager {
  final _controller = StreamController<CoreEvent>();

  CoreEventManager._() {
    _controller.stream.listen((event) {
      switch (event.type) {
        case CoreEventType.log:
          final log = Log.fromJson(event.data);
          for (final CoreEventListener listener in _listeners) {
            listener.onLog(log);
          }
          break;
        case CoreEventType.delay:
          final delay = Delay.fromJson(event.data);
          for (final CoreEventListener listener in _listeners) {
            listener.onDelay(delay);
          }
          break;
        case CoreEventType.request:
          final tracker = TrackerInfo.fromJson(event.data);
          for (final CoreEventListener listener in _listeners) {
            listener.onRequest(tracker);
          }
          break;
        case CoreEventType.loaded:
          for (final CoreEventListener listener in _listeners) {
            listener.onLoaded(event.data);
          }
          break;
        case CoreEventType.crash:
          for (final CoreEventListener listener in _listeners) {
            listener.onCrash(event.data);
          }
          break;
      }
    });
  }

  static final CoreEventManager instance = CoreEventManager._();

  final ObserverList<CoreEventListener> _listeners =
      ObserverList<CoreEventListener>();

  bool get hasListeners {
    return _listeners.isNotEmpty;
  }

  void sendEvent(CoreEvent event) {
    _controller.add(event);
  }

  void addListener(CoreEventListener listener) {
    _listeners.add(listener);
  }

  void removeListener(CoreEventListener listener) {
    _listeners.remove(listener);
  }
}

final coreEventManager = CoreEventManager.instance;
