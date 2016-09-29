type AsyncCallback = (err: ?Object, results: ?any) => void;
type ParallelFuncs = (callback: AsyncCallback) => void;

declare module 'async' {
  declare function parallel(tasks: Array<ParallelFuncs> | Object,
                            callback: AsyncCallback): void
  declare function series(tasks: Array<ParallelFuncs>,
                          callback: AsyncCallback): void
  declare function each(arr: Array<any>,
                        func: Function,
                        callback: Function): void
}