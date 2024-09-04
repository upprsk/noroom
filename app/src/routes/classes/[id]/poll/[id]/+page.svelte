<script lang="ts">
  import { zPollAnswerWithUserSchema, zPollSchema } from '$lib/models';
  import { pb, updateFromEvent, updateOneFromEvent } from '$lib/pocketbase';
  import { onMount } from 'svelte';
  import { z } from 'zod';
  import PollAnswers from './PollAnswers.svelte';
  import PollEdit from './PollEdit.svelte';
  import PollQuestion from './PollQuestion.svelte';

  type Poll = z.infer<typeof zPollSchema>;
  type PollAnswer = z.infer<typeof zPollAnswerWithUserSchema>;

  export let data;

  let pollUnavailable = false;

  let unsubPoll: () => void;
  const subPoll = async () => {
    unsubPoll = await pb.collection('polls').subscribe<Poll>(data.poll.id, (e) => {
      console.log(e);
      data.poll = updateOneFromEvent(e, zPollSchema, (v) => {
        pollUnavailable = true;
        return v;
      });
    });
  };

  let unsubAnswers: () => void;
  const subAnswers = async () => {
    unsubAnswers = await pb
      .collection('pollAnswers')
      .subscribe<PollAnswer>(
        '*',
        (e) => (data.answers = updateFromEvent(e, zPollAnswerWithUserSchema, data.answers)),
        { filter: pb.filter('poll={:poll}', { poll: data.poll.id }), expand: 'user' },
      );
  };

  onMount(() => {
    subPoll();
    subAnswers();

    return () => {
      unsubPoll();
      unsubAnswers();
    };
  });
</script>

{#if pollUnavailable}
  <div role="alert" class="alert alert-warning">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-6 w-6 shrink-0 stroke-current"
      fill="none"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
      />
    </svg>
    <span>
      Esta poll foi deletada no sistema. <a href="/classes/{data.poll.class}" class="link">voltar</a
      >
    </span>
  </div>
{/if}

{#if data.user.role === 'editor'}
  <PollAnswers user={data.user} klass={data.klass} poll={data.poll} answers={data.answers} />
  <PollEdit poll={data.poll} />
  {#if data.poll.active && !data.answers.find((it) => it.user === data.user.id)}
    <PollQuestion user={data.user} poll={data.poll} />
  {/if}
{:else if data.answers.find((it) => it.user === data.user.id)}
  <!-- user has already answered, show results -->
  <PollAnswers user={data.user} klass={data.klass} poll={data.poll} answers={data.answers} />
{:else}
  <!-- user has not answered, show question -->
  <PollQuestion user={data.user} poll={data.poll} />
{/if}
