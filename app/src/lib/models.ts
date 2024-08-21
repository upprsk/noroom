import { z } from 'zod';

export const zModelBase = z.object({
  id: z.string(),
  created: z.string(),
  updated: z.string(),
});

export const zErrorDataItemSchema = z.object({
  code: z.string(),
  message: z.string(),
});

export const zErrorSchema = z.object({
  code: z.number().int(),
  message: z.string(),
});

export const zMakeErrorDataSchema = <T extends z.ZodTypeAny>(keys: T) =>
  zErrorSchema.extend({
    data: z.record(keys, zErrorDataItemSchema),
  });

export const zUserSchema = zModelBase.extend({
  username: z.string(),
  email: z.string().email(),
  emailVisibility: z.boolean(),
  verified: z.boolean(),
  name: z.string(),
  mat: z.number().int(),
  curso: z.string(),
  avatar: z.string(),
});

export const zRegisterSchema = zUserSchema
  .pick({
    username: true,
    email: true,
    name: true,
    mat: true,
    curso: true,
  })
  .extend({
    password: z.string(),
    passwordConfirm: z.string(),
  });

export const zLogin = z.object({
  username: z.string(),
  password: z.string(),
});
