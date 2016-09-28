export const INTERNAL_ROUTER_READY = 'INTERNAL_ROUTER_READY';
export const INTERNAL_STORE_IDLE = 'INTERNAL_STORE_IDLE';

export function internalRouterReady() {
  return {
    type: INTERNAL_ROUTER_READY,
  };
}

export function internalStoreIdle(state) {
  return {
    type: INTERNAL_STORE_IDLE,
    payload: state,
  };
}
