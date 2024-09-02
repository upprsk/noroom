import { zLogin, zMakeErrorDataSchema } from '$lib/models';

export const zFormSchema = zLogin;
export const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
