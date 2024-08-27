<script lang="ts">
  import { invalidate } from '$app/navigation';
  import BasicCard from '$lib/components/BasicCard.svelte';
  import { zPodSchema, zPodServerSchema, zPodServerWithPodsSchema } from '$lib/models.js';
  import { pb, updateFromEvent } from '$lib/pocketbase.js';
  import { ClientResponseError } from 'pocketbase';
  import { onMount } from 'svelte';
  import type { z } from 'zod';

  export let data;

  let isStarting = false;
  const startPod = async (id: string) => {
    try {
      isStarting = true;

      const res = await pb.send(`/api/noroom/pod/${id}/start`, {
        method: 'post',
      });

      console.log(res);
    } catch (e) {
      console.error(e);

      if (e instanceof ClientResponseError) {
        console.error(e.message);
      }
    } finally {
      isStarting = false;
    }
  };

  let isStopping = false;
  const stopPod = async (id: string) => {
    try {
      isStopping = true;

      const res = await pb.send(`/api/noroom/pod/${id}/stop`, {
        method: 'post',
      });

      console.log(res);
    } catch (e) {
      console.error(e);

      if (e instanceof ClientResponseError) {
        console.error(e.message);
      }
    } finally {
      isStopping = false;
    }
  };

  let isKilling = false;
  const killPod = async (id: string) => {
    try {
      isKilling = true;

      const res = await pb.send(`/api/noroom/pod/${id}/kill`, {
        method: 'post',
      });

      console.log(res);
    } catch (e) {
      console.error(e);

      if (e instanceof ClientResponseError) {
        console.error(e.message);
      }
    } finally {
      isKilling = false;
    }
  };

  let isInspecting = false;
  const inspectPod = async (id: string) => {
    try {
      isInspecting = true;

      const res = await pb.send(`/api/noroom/pod/${id}/inspect`, {
        method: 'post',
      });

      console.log(res);
      invalidate('app:podServers');
    } catch (e) {
      console.error(e);

      if (e instanceof ClientResponseError) {
        console.error(e.message);
      }
    } finally {
      isInspecting = false;
    }
  };

  // --------------------------------------------------------------------------

  type PodServer = z.infer<typeof zPodServerWithPodsSchema>;
  type Pod = z.infer<typeof zPodSchema>;

  let unsubPodServers: () => void;
  const subPodServers = async () => {
    unsubPodServers = await pb
      .collection('podServers')
      .subscribe<PodServer>(
        '*',
        (e) => (data.podServers = updateFromEvent(e, zPodServerSchema, data.podServers)),
      );
  };

  // --------------------------------------------------------------------------

  let unsubPods: () => void;
  const subPods = async () => {
    unsubPods = await pb
      .collection('pods')
      .subscribe<Pod>('*', (e) => (data.pods = updateFromEvent(e, zPodSchema, data.pods)));
  };

  // --------------------------------------------------------------------------

  onMount(() => {
    subPodServers();
    subPods();

    return () => {
      if (unsubPodServers) unsubPodServers();
      if (unsubPods) unsubPods();
    };
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">Servidores Pod</svelte:fragment>

  <div class="h-5"></div>

  <div class="prose">
    <ul>
      {#each data.podServers as srv (srv.id)}
        {@const pods = data.pods.filter((p) => p.server === srv.id)}
        <li>
          <span><b>{srv.name}:</b> {srv.address}</span>
          <ul>
            {#each pods as pod (pod.id)}
              <li>
                <div class="flex justify-between">
                  <span>
                    <b>{pod.name} - {pod.image}:</b>
                    {pod.status}
                    {#if pod.running}
                      <span class="badge badge-success">running</span>
                    {:else}
                      <span class="badge badge-ghost">stopped</span>
                    {/if}
                  </span>

                  <div>
                    <button
                      type="button"
                      class="btn btn-sm"
                      disabled={isStarting}
                      on:click={() => startPod(pod.id)}
                    >
                      start
                    </button>
                    <button
                      type="button"
                      class="btn btn-sm"
                      disabled={isStopping}
                      on:click={() => stopPod(pod.id)}
                    >
                      stop
                    </button>
                    <button
                      type="button"
                      class="btn btn-info btn-sm"
                      disabled={isInspecting}
                      on:click={() => inspectPod(pod.id)}
                    >
                      inspect
                    </button>
                    <button
                      type="button"
                      class="btn btn-warning btn-sm"
                      disabled={isKilling}
                      on:click={() => killPod(pod.id)}
                    >
                      kill
                    </button>
                  </div>
                </div>
              </li>
            {/each}
          </ul>
        </li>
      {/each}
    </ul>
  </div>
</BasicCard>
