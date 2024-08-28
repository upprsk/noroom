<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { zPodWithServerSchema } from '$lib/models.js';
  import { pb, simpleSend, updateOneFromEvent } from '$lib/pocketbase.js';
  import '@xterm/xterm/css/xterm.css';
  import { onMount } from 'svelte';
  import { z } from 'zod';
  import Terminal from './Terminal.svelte';

  export let data;

  const {
    loading: loadingStart,
    errors: errorsStart,
    send: sendStart,
  } = simpleSend(pb, z.unknown(), `/api/noroom/pod/${data.pod.id}/start`, {
    method: 'post',
  });

  const {
    loading: loadingStop,
    errors: errorsStop,
    send: sendStop,
  } = simpleSend(pb, z.unknown(), `/api/noroom/pod/${data.pod.id}/stop`, {
    method: 'post',
  });

  const {
    loading: loadingKill,
    errors: errorsKill,
    send: sendKill,
  } = simpleSend(pb, z.unknown(), `/api/noroom/pod/${data.pod.id}/kill`, {
    method: 'post',
  });

  const {
    loading: loadingInspect,
    errors: errorsInspect,
    send: sendInspect,
  } = simpleSend(pb, z.unknown(), `/api/noroom/pod/${data.pod.id}/inspect`, {
    method: 'post',
  });

  $: anyLoading = $loadingStart || $loadingStop || $loadingKill || $loadingInspect;

  // --------------------------------------------------------------------------

  let wantsConnectedToPod = false;

  // --------------------------------------------------------------------------

  type Pod = z.infer<typeof zPodWithServerSchema>;

  let podUnavailable = false;

  let unsubPods: () => void;
  const subPods = async () => {
    unsubPods = await pb.collection('pods').subscribe<Pod>(
      data.pod.id,
      (e) =>
        (data.pod = updateOneFromEvent(e, zPodWithServerSchema, (v) => {
          podUnavailable = true;
          return v;
        })),
      { expand: 'server' },
    );
  };

  // --------------------------------------------------------------------------

  onMount(() => {
    subPods();

    return () => {
      if (unsubPods) unsubPods();
    };
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">
    {data.pod.name} - {data.pod.image}

    {#if data.pod.running}
      <span class="badge badge-success">running</span>
    {/if}
  </svelte:fragment>

  {#if data.pod.expand?.server}
    <h5 class="text-slate-400">{data.pod.expand.server.name} server</h5>
  {/if}

  {#if podUnavailable}
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
      <span>Este pod foi deletado no sistema. <a href="/pods" class="link">voltar</a></span>
    </div>
  {/if}

  <div class="h-5"></div>

  <div class="flex w-full items-center gap-2">
    <button type="button" class="btn btn-info btn-sm" disabled={anyLoading} on:click={sendInspect}>
      inspect
    </button>

    <button type="button" class="btn btn-primary btn-sm" disabled={anyLoading} on:click={sendStart}>
      start
    </button>

    <button type="button" class="btn btn-accent btn-sm" disabled={anyLoading} on:click={sendStop}>
      stop
    </button>

    <button type="button" class="btn btn-warning btn-sm" disabled={anyLoading} on:click={sendKill}>
      kill
    </button>

    <div class="form-control">
      <label class="label cursor-pointer gap-2">
        <span class="label-text">Connected</span>
        <input
          type="checkbox"
          disabled={!data.pod.running}
          bind:checked={wantsConnectedToPod}
          class="checkbox"
        />
      </label>
    </div>
  </div>

  <ErrorAlert errors={$errorsStart} />
  <ErrorAlert errors={$errorsStop} />
  <ErrorAlert errors={$errorsKill} />
  <ErrorAlert errors={$errorsInspect} />

  <div class="h-5"></div>

  <Terminal bind:wantsConnectedToPod pod={data.pod} />
</BasicCard>
