<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import EmailInput from '$lib/components/input/EmailInput.svelte';
  import IntInput from '$lib/components/input/IntInput.svelte';
  import PasswordInput from '$lib/components/input/PasswordInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { superForm } from 'sveltekit-superforms';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(data.form);
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
>
  <svelte:fragment slot="title">Sign-Up</svelte:fragment>

  <TextInput
    name="name"
    errors={$errors.name}
    bind:value={$form.name}
    constraints={$constraints.name}
  >
    Nome
  </TextInput>

  <EmailInput
    name="email"
    errors={$errors.email}
    bind:value={$form.email}
    constraints={$constraints.email}
  >
    Email
  </EmailInput>

  <TextInput
    name="username"
    errors={$errors.username}
    bind:value={$form.username}
    constraints={$constraints.username}
  >
    Username
  </TextInput>

  <IntInput name="mat" errors={$errors.mat} bind:value={$form.mat} constraints={$constraints.mat}>
    Matrícula
  </IntInput>

  <TextInput
    name="curso"
    errors={$errors.curso}
    bind:value={$form.curso}
    constraints={$constraints.curso}
  >
    Curso
  </TextInput>

  <PasswordInput
    name="password"
    errors={$errors.password}
    bind:value={$form.password}
    constraints={$constraints.password}
  >
    Senha
  </PasswordInput>

  <PasswordInput
    name="passwordConfirm"
    errors={$errors.passwordConfirm}
    bind:value={$form.passwordConfirm}
    constraints={$constraints.passwordConfirm}
  >
    Confirme a Senha
  </PasswordInput>

  <svelte:fragment slot="notes">
    <p>Já possui uma conta? <a href="/login" class="link">Login</a></p>
  </svelte:fragment>
</BasicFormCard>
