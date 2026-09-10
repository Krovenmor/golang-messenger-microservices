const MAX_AVATAR_SIDE = 600
const JPEG_QUALITY = 0.92
const AVATAR_FILE_NAME = 'avatar.jpg'
const AVATAR_MIME = 'image/jpeg'

export class UnreadableImageError extends Error {
  constructor() {
    super('image could not be decoded by the browser')
  }
}

export class AvatarBudgetExceededError extends Error {
  constructor(
    readonly encodedBytes: number,
    readonly budgetBytes: number
  ) {
    super('encoded avatar does not fit into the remaining media space')
  }
}

function loadImage(file: File): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(file)
    const img = new Image()
    img.onload = () => {
      URL.revokeObjectURL(url)
      resolve(img)
    }
    img.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new UnreadableImageError())
    }
    img.src = url
  })
}

function drawCenteredSquare(img: HTMLImageElement, targetSide: number): HTMLCanvasElement {
  const sourceSide = Math.min(img.width, img.height)
  const sourceX = (img.width - sourceSide) / 2
  const sourceY = (img.height - sourceSide) / 2

  const canvas = document.createElement('canvas')
  canvas.width = targetSide
  canvas.height = targetSide
  const ctx = canvas.getContext('2d')
  if (!ctx) throw new Error('Canvas недоступен')
  ctx.drawImage(img, sourceX, sourceY, sourceSide, sourceSide, 0, 0, targetSide, targetSide)
  return canvas
}

function encodeJpeg(canvas: HTMLCanvasElement): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob((blob) => (blob ? resolve(blob) : reject(new UnreadableImageError())), AVATAR_MIME, JPEG_QUALITY)
  })
}

export async function prepareAvatarFile(file: File, budgetBytes = Number.POSITIVE_INFINITY): Promise<File> {
  const img = await loadImage(file)
  if (!img.width || !img.height) throw new UnreadableImageError()

  const targetSide = Math.min(img.width, img.height, MAX_AVATAR_SIDE)
  const blob = await encodeJpeg(drawCenteredSquare(img, targetSide))

  if (blob.size > budgetBytes) throw new AvatarBudgetExceededError(blob.size, budgetBytes)

  return new File([blob], AVATAR_FILE_NAME, { type: AVATAR_MIME })
}
