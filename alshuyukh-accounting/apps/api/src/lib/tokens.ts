import { createHash, randomBytes } from 'node:crypto';
import { SignJWT, jwtVerify } from 'jose';
import { unauthorized } from './errors.js';

export interface AccessClaims {
  sub: string; // user id
  tid: string; // active tenant id
  sid: string; // session id
}

export class TokenService {
  private readonly key: Uint8Array;

  constructor(secret: string, private readonly accessTtlSeconds: number) {
    this.key = new TextEncoder().encode(secret);
  }

  signAccess(claims: AccessClaims): Promise<string> {
    return new SignJWT({ tid: claims.tid, sid: claims.sid, typ: 'access' })
      .setProtectedHeader({ alg: 'HS256' })
      .setSubject(claims.sub)
      .setIssuer('alshuyukh')
      .setAudience('alshuyukh-api')
      .setIssuedAt()
      .setExpirationTime(`${this.accessTtlSeconds}s`)
      .sign(this.key);
  }

  async verifyAccess(token: string): Promise<AccessClaims> {
    try {
      const { payload } = await jwtVerify(token, this.key, {
        issuer: 'alshuyukh',
        audience: 'alshuyukh-api',
        algorithms: ['HS256'],
      });
      if (payload.typ !== 'access' || typeof payload.sub !== 'string' ||
          typeof payload.tid !== 'string' || typeof payload.sid !== 'string') {
        throw new Error('bad claims');
      }
      return { sub: payload.sub, tid: payload.tid, sid: payload.sid };
    } catch {
      throw unauthorized('Invalid or expired token');
    }
  }
}

export const newRefreshToken = () => randomBytes(32).toString('base64url');
export const hashToken = (token: string) => createHash('sha256').update(token).digest('hex');
