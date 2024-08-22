import { PUBLIC_POCKETBASE_URL } from '$env/static/public';
import PocketBase, {
  ClientResponseError,
  type BaseModel,
  type FileOptions,
  type RecordSubscription,
} from 'pocketbase';
import { setError, type Infer, type SuperValidated } from 'sveltekit-superforms';
import type { z } from 'zod';

export function createInstance() {
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
