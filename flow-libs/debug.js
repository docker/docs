type DebugFunction = (thing: any) => void;

declare module 'debug' {
  declare function exports(string: string): DebugFunction;
}