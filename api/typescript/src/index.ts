import Cache from "./modules/cache";
import GraphQL from "./modules/graphql";
import Http from "./modules/http";
import Log from "./modules/log";
import Request from "./modules/request";
import Runnable from "./modules/runnable";

import FFI from "./modules/ffi";

export const cache: Cache = new Cache();
export const graphql: GraphQL = new GraphQL();
export const http: Http = new Http();
export const log: Log = new Log();
export const request: Request = new Request();
export const runnable: Runnable = new Runnable();

export function setup(imports: object, ident: number) {
  // Apply the import object
  // @ts-ignore
  FFI.env._exports = imports;

  FFI.ident = ident;
}
