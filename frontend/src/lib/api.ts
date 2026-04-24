export type User = { id: number; email: string }

export async function api<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const response = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers ?? {}),
    },
    credentials: 'include',
  })

  if (!response.ok) {
    let message = `HTTP ${response.status}`
    try {
      const data = await response.json()
      if (data?.error) message = data.error
    } catch {
      /* empty */
    }
    throw new Error(message)
  }
  const text = await response.text()
  return (text ? JSON.parse(text) : null) as T
}
