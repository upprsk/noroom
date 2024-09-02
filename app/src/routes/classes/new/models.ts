import { zClassSchema, zMakeErrorDataSchema } from '$lib/models';

export const zFormSchema = zClassSchema;
export const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
