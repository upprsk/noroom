import { zClassSchema, zPollAnswerWithUserSchema, zPollSchema } from '$lib/models';
import { pb } from '$lib/pocketbase';
import { currentUser } from '$lib/stores/user';
import { type Load, redirect } from '@sveltejs/kit';
import { get } from 'svelte/store';

const zPollAnswerArraySchema = zPollAnswerWithUserSchema.array();

export const load: Load = async ({ params, fetch }) => {
  const user = get(currentUser);
  const { id } = params;

  if (!id || !user) throw redirect(303, '/');

  const pollP = pb
    .collection('polls')
    .getOne(id, { fetch })
    .then((r) => zPollSchema.parse(r));

  const answersP = pb
    .collection('pollAnswers')
    .getFullList({
      fetch,
      filter: pb.filter('poll={:poll}', { poll: id }),
      expand: 'user',
    })
    .then((r) => zPollAnswerArraySchema.parse(r));

  const [poll, answers] = await Promise.all([pollP, answersP]);

  const klass = await pb
    .collection('classes')
    .getOne(poll.class, { fetch })
    .then((r) => zClassSchema.parse(r));

  return { user, klass, poll, answers };
};
