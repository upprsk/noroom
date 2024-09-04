<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import type {
    zClassSchema,
    zPollAnswerWithUserSchema,
    zPollSchema,
    zUserSchema,
  } from '$lib/models';
  import type { z } from 'zod';
  import PieChart from './PieChart.svelte';

  export let user: z.infer<typeof zUserSchema>;
  export let klass: z.infer<typeof zClassSchema>;
  export let poll: z.infer<typeof zPollSchema>;
  export let answers: z.infer<typeof zPollAnswerWithUserSchema>[];

  $: data = Object.entries(
    answers
      .map((a) => a.answer)
      .reduce(
        (acc, v) => {
          if (Array.isArray(v)) {
            for (const it of v) acc[it] = (acc[it] ?? 0) + 1;
          } else {
            acc[v] = (acc[v] ?? 0) + 1;
          }

          return acc;
        },
        {} as Record<string | number, number>,
      ),
  ).map(([name, value]) => ({ name, value }));
</script>

<BasicCard>
  <svelte:fragment slot="title">Resultados</svelte:fragment>

  {#if !poll.active}
    <div role="alert" class="alert">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="h-6 w-6 shrink-0 stroke-info"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
      <span>Esta poll foi encerrada</span>
    </div>
  {/if}

  {#if poll.expects === 'option' || poll.expects === 'multi'}
    <PieChart title={poll.title} subtitle={klass.title} {data} />
  {:else if poll.expects === 'number' || poll.expects === 'string'}
    <div class="overflow-x-auto">
      <table class="table">
        <!-- head -->
        <thead>
          <tr>
            <th></th>
            <th>Nome</th>
            <th>Resposta</th>
          </tr>
        </thead>
        <tbody>
          {#each answers as answer, idx (answer.id)}
            <tr class:bg-base-200={answer.user === user.id}>
              <th>{idx + 1}</th>
              <td>{answer.expand?.user.name}</td>
              <td>{answer.answer}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else}
    <ErrorAlert errors={['Something is wrong']} />
  {/if}
</BasicCard>
