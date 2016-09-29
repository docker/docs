export type FluxibleActionContext = {
  dispatch(eventName: string,
           payload: any): void;
}
