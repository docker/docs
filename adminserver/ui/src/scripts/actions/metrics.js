'use strict';

import { Metrics } from 'dtr-js-sdk';
import Rx from 'rx';

import consts from 'consts';
const { metrics } = consts;
// if ("production" !== process.env.NODE_ENV) {
//   require("./__tests__/mocks/metrics");
// }

/**
 * Long polling interval in milliseconds
 *
 */
const INTERVAL = 2000;

/**
 * Stores the observable from polling metrics and the subscription we use to
 * listen to said observable.
 *
 */
let observable, subscription;

/**
 * This stores the last lookup time and CPU wall usage of the last fetched data,
 * which allows us to calculate CPU usage of the next lookup using wall time
 * information.
 *
 * previousCPU is an object as such:
 * {
 *   global: int,
 *   containers: {
 *     containerid1: int,
 *     containerid2: int,
 *     ...
 *   }
 * }
 */
let previousCPU = {}, previousNetwork = {}, previousTime = 0;

function unsubscribe() {
  if (subscription) {
    subscription.dispose();
    subscription = null;
  }
}

function filterMetrics(response) {

  // Normalise our metrics
  return response.data.map( raw => {
    if (raw.time === previousTime) {
      return false;
    }

    // Normalise global hard disk values into bytes (from KB)
    let { total, available } = raw.disks['/'];
    total *= 1024;
    available *= 1024;
    const used = total - available;

    // Get the time difference between the previous request and now for
    // capturing CPU information
    const timeDiffMs = new Date(raw.time) - new Date(previousTime);
    const timeDiffNs = timeDiffMs * Math.pow(10, 6);
    const timeDiffS = timeDiffMs * Math.pow(10, -3);

    // By default this should be 0 percent if we have no previous
    // information to calculate percentage from
    var globalCpu = 0;
    if (previousTime) {
      // CPU stats are defined in terms of nanoseconds of CPU time so use timeDiff in ns
      // Then normalise global CPU information into a percentage
      globalCpu = 100 * (raw.total_cpu - previousCPU.global) / timeDiffNs;
    }
    if (isNaN(globalCpu)) {
      globalCpu = 0;
    }
    previousCPU.global = raw.total_cpu;

    // Normalise each container's data
    let containers = {};
    Object.keys(raw.metrics).forEach( containerId => {
      const rawMetrics = raw.metrics[containerId];
      const keyMap = {
        tx_bytes: 'tx',
        rx_bytes: 'rx',
        tx_dropped: 'txDropped',
        rx_dropped: 'rxDropped'
      };

      let cpu = 0;
      let networkStats = {};
      Object.keys(keyMap).forEach(statKey => networkStats[keyMap[statKey]] = 0);
      if (previousTime) {
        if (previousCPU[containerId]) {
          // CPU stats are defined in terms of nanoseconds of CPU time so use timeDiff in ns
          // Then normalise global CPU information into a percentage
          cpu = 100 * (rawMetrics.cpu.total - previousCPU[containerId]) / timeDiffNs;
        }
        if (previousNetwork[containerId]) {
          Object.keys(keyMap).forEach(statKey => {
            networkStats[keyMap[statKey]] = (rawMetrics.network_stats[statKey] - previousNetwork[containerId][statKey]) / timeDiffS;
          });
        }
      }
      if (isNaN(cpu)) {
        cpu = 0;
      }

      // Update previous CPU information
      previousCPU[containerId] = rawMetrics.cpu.total;
      previousNetwork[containerId] = rawMetrics.network_stats;

      const containerName = window.containerIDs[containerId];
      containers[containerName] = {
        cpu: {
          percentage: cpu || 0
        },
        memory: {
          usage: rawMetrics.memory.usage,
          maxUsage: rawMetrics.memory.max_usage,
          limit: rawMetrics.memory.limit
        },
        network: networkStats
      };
    });

    // After calculations we can set previous time to the current response
    previousTime = raw.time;

    return {
      time: raw.time,
      raw,
      filtered: {
        cpu: {
          percentage: globalCpu
        },
        storage: {
          total,
          available,
          used
        },
        containers
      }
    };
  }).filter(Boolean);
}

export function sendClientAnalytics() {
  return {
    type: metrics.SEND_ANALYTICS,
    meta: {
      promise: Metrics.sendClientAnalytics()
    }
  };
}

export function getMetrics() {
  return {
    type: metrics.FETCH_METRICS,
    meta: {
      promise: Metrics.getMetrics().then(response => filterMetrics(response))
    }
  };
}

export function observeMetrics() {
  return dispatch => {
    unsubscribe();
    observable = Rx.Observable
      .interval(INTERVAL)
      .flatMap(() => Rx.Observable.fromPromise(Metrics.getMetrics(new Date(previousTime))))
      .retry(60)
      .map(filterMetrics);
    subscription = observable.subscribe(
      data => { dispatch({ type: metrics.OBSERVE_METRICS_DATA, payload: data }); },
    );
    dispatch({ type: metrics.OBSERVE_METRICS });
  };
}

export function unobserveMetrics() {
  unsubscribe();
  return { type: metrics.UNOBSERVE_METRICS };
}
