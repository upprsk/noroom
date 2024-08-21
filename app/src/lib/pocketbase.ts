import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import PocketBase, { ClientResponseError, type BaseModel, type FileOptions } from 'pocketbase';
import { setError, type Infer, type SuperValidated } from 'sveltekit-superforms';
import type { z } from 'zod';

export function createInstance() {
  return new PocketBase(PUBLIC_POCKETBASE_URL);
}

export const pb = createInstance();

export const processError = <T extends z.ZodTypeAny, S extends z.ZodTypeAny>(
  form: SuperValidated<Infer<T>>,
  e: unknown,
  schema: S,
) => {
  if (e instanceof ClientResponseError) {
    const { success, data, error } = schema.safeParse(e.data);
    if (!success) {
      console.error('success is false', error.issues);
      return setError(form, e.message);
    }

    for (const [k, v] of Object.entries(data.data)) {
      const { message } = v as { message: string };
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      setError(form, k, message);
    }

    return setError(form, e.message);
  }
};

export const getImageUrl = (m: BaseModel, file: string, opt?: FileOptions) =>
  pb.files.getUrl(m, file, opt);

export const sendFingerprint = async (userid: string | undefined, fingerprint: string) => {
  try {
    await pb.send('/api/noroom/tracking', {
      method: 'POST',
      body: { userid, fingerprint },
    });
  } catch (e) {
    console.error(e);
  }
};
