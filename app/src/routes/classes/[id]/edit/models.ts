import { zClassSchema, zMakeErrorDataSchema } from '$lib/models';

export const zFormSchema = zClassSchema;
export const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
export const zSaveSchema = zClassSchema.omit({ attachments: true });
