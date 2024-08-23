import { z } from 'zod';

export const zModelBase = z.object({
  id: z.string(),
  created: z.string(),
  updated: z.string(),
  collectionId: z.string(),
  collectionName: z.string(),
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
  email: z.string().email().optional(),
  emailVisibility: z.boolean(),
  verified: z.boolean(),
  name: z.string(),
  mat: z.number().int(),
  curso: z.string(),
  avatar: z.string(),
  role: z.enum(['editor', 'student']),
});

export const zEndDeviceSchema = zModelBase.extend({
  fingerprint: z.string(),
  owners: z.string().array(),
  deviceData: z
    .object({
      locales: z.object({
        languages: z.string(),
        timezone: z.string(),
      }),
      system: z.object({
        platform: z.string(),
        product: z.string(),
        productSub: z.string(),
        browser: z.object({ name: z.string(), version: z.string() }),
      }),
      hardware: z.object({
        architecture: z.number(),
        deviceMemory: z.string(),
        jsHeapSizeLimit: z.number(),
      }),
    })
    .passthrough()
    .nullish(),
  locationData: z
    .object({
      query: z.string(),
      status: z.string(),
      country: z.string(),
      countryCode: z.string(),
      region: z.string(),
      regionName: z.string(),
      city: z.string(),
      zip: z.string(),
      lat: z.number(),
      lon: z.number(),
      timezone: z.string(),
      isp: z.string(),
      org: z.string(),
      as: z.string(),
    })
    .passthrough()
    .nullish(),
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

export const zClassPresence = zModelBase.extend({
  class: z.string(),
  fingerprint: z.string(),
  user: z.string(),
});

export const zClassPresenceWithUser = zClassPresence.extend({
  expand: z.object({
    user: zUserSchema,
  }),
});

export const zClassSchema = zModelBase.extend({
  title: z.string(),
  content: z.string(),
  attachments: z.string().array(),
  owner: z.string(),
  live: z.boolean(),
  latitude: z.number(),
  longitude: z.number(),
  radius: z.number(),
});

export const zClassWithPresenceSchema = zClassSchema.extend({
  expand: z
    .object({
      classPresenceEntries_via_class: zClassPresenceWithUser.array(),
    })
    .optional(),
});

export const zFileUploadSchema = z.object({
  attachments: z
    .instanceof(File, { message: 'Select a file' })
    .refine((f) => f.size < 5_242_880, 'Max 5MiB upload size')
    .array(),
});

export type BaseModel = z.infer<typeof zModelBase>;
export type UserModel = z.infer<typeof zUserSchema>;
export type ClassModel = z.infer<typeof zClassSchema>;
