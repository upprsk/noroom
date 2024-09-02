import { zRegisterSchema, zMakeErrorDataSchema } from '$lib/models';

export const zFormSchema = zRegisterSchema;
export const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
