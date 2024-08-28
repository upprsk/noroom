import { browser } from '$app/environment';
import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import PocketBase, {
  ClientResponseError,
  type BaseModel,
  type FileOptions,
  type RecordSubscription,
  type SendOptions,
} from 'pocketbase';
import { derived, get, writable } from 'svelte/store';
import { setError, type Infer, type SuperValidated } from 'sveltekit-superforms';
import type { z } from 'zod';

export function createInstance() {
  if (browser) return new PocketBase('/');
  return new PocketBase(PUBLIC_POCKETBASE_URL);
}

export const pb = createInstance();

export const updateFromEvent = <T extends z.ZodTypeAny, U extends z.infer<T>>(
  e: RecordSubscription<U>,
  schema: T,
  data: U[],
): U[] => {
  const del = () => data.filter((item) => item.id !== e.record.id);
  const create = () => [...data, e.record];
  const update = () => {
    const idx = data.findIndex((item) => item.id == e.record.id);
    if (idx >= 0) {
      data[idx] = schema.parse(e.record);
    }

    return data;
  };

  switch (e.action) {
    case 'update':
      return update();
    case 'create':
      return create();
    case 'delete':
      return del();
    default:
      throw new Error(`invalid event action: ${e.action}`);
  }
};

export const updateOneFromEvent = <T extends z.ZodTypeAny, U extends z.infer<T>>(
  e: RecordSubscription<U>,
  schema: T,
  onDelete: (v: U) => U,
): U => {
  const record = schema.parse(e.record);

  switch (e.action) {
    case 'update':
      return schema.parse(record);
    case 'delete':
      return onDelete(record);
    default:
      throw new Error(`invalid event action: ${e.action}`);
  }
};

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

export const getFileUrl = (m: BaseModel, file: string, opt?: FileOptions) =>
  pb.files.getUrl(m, file, opt);

export const simpleSend = <T extends z.ZodTypeAny>(
  pb: PocketBase,
  schema: T,
  path: string,
  options: SendOptions,
) => {
  const loading = writable(false);
  const error = writable<ClientResponseError | undefined>();
  const data = writable<z.infer<T>>();

  const errors = derived(error, ($error) => {
    if (!$error) return undefined;

    return [$error.message];
  });

  return {
    loading,
    error,
    errors,
    data,
    send: async () => {
      if (get(loading)) throw new Error('already sending');

      try {
        loading.set(true);
        error.set(undefined);

        const res = schema.parse(await pb.send(path, options)) as z.infer<T>;
        data.set(res);
      } catch (e) {
        console.error('error in send to ', path, ':', e);

        if (e instanceof ClientResponseError) {
          error.set(e);
        }
      } finally {
        loading.set(false);
      }
    },
  };
};
