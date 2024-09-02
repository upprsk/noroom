<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import PasswordInput from '$lib/components/input/PasswordInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { defaults, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import { zErrorSchema, zFormSchema } from './models';
  import { pb, processError } from '$lib/pocketbase';
  import { goto } from '$app/navigation';

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    defaults(zod(zFormSchema)),
    {
      SPA: true,
      validators: zod(zFormSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        try {
          await pb.collection('users').authWithPassword(form.data.username, form.data.password);
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }

        goto('/');
      },
    },
  );
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
>
  <svelte:fragment slot="title">Sign-In</svelte:fragment>

  <TextInput
    name="username"
    errors={$errors.username}
    bind:value={$form.username}
    constraints={$constraints.username}
  >
    Username or Email
  </TextInput>

  <PasswordInput
    name="password"
    errors={$errors.password}
    bind:value={$form.password}
    constraints={$constraints.password}
  >
    Senha
  </PasswordInput>

  <svelte:fragment slot="notes">
    <p>Ainda não possui uma conta? <a href="/register" class="link">Crie aqui!</a></p>
  </svelte:fragment>
</BasicFormCard>
