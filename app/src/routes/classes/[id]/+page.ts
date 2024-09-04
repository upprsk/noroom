import { zClassWithPresenceSchema, zPollAnswerSchema, zPollSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { redirect, type Load } from '@sveltejs/kit';
import { get } from 'svelte/store';
import { z } from 'zod';

const zPollWithAnswersSchema = zPollSchema.extend({
  expand: z
    .object({
      pollAnswers_via_poll: zPollAnswerSchema.array(),
    })
    .optional()
    .default({
      pollAnswers_via_poll: [],
    }),
});

const zPollArraySchema = zPollWithAnswersSchema.array();

export const load: Load = async ({ params, fetch }) => {
  const { id } = params;
  const user = get(currentUser);

  if (!id || !user) throw redirect(303, '/');

  const classP = pb
    .collection('classes')
    .getOne(id, {
      fetch,
      expand: 'classPresenceEntries_via_class,classPresenceEntries_via_class.user',
    })
    .then((r) => zClassWithPresenceSchema.parse(r));

  const classPollsP = pb
    .collection('polls')
    .getFullList({
      fetch,
      filter: pb.filter('class={:klass}', { klass: id }),
      expand: 'pollAnswers_via_poll',
    })
    .then((r) => zPollArraySchema.parse(r));

  const [klass, classPolls] = await Promise.all([classP, classPollsP]);

  return { user, klass, classPolls };
};
