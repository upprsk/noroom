import { zMakeErrorDataSchema, zPodSchema, zPodServerSchema } from '$lib/models';

export const zFormSchema = zPodSchema;
export const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());
export const zPodServerArraySchema = zPodServerSchema.array();
