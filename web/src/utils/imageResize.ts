// The backend caps avatars at 600x600px (png/jpeg). Rather than reject
// whatever photo the person picks, crop it to a centered square and scale
// it down to fit — so any reasonably-sized photo just works.

const MAX_AVATAR_SIDE = 600

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
      reject(new Error('Не удалось прочитать изображение'))
    }
    img.src = url
  })
}

export async function prepareAvatarFile(file: File, maxSide = MAX_AVATAR_SIDE): Promise<File> {
  const img = await loadImage(file)
  const side = Math.min(img.width, img.height)
  const sx = (img.width - side) / 2
  const sy = (img.height - side) / 2
  const targetSide = Math.min(side, maxSide)

  const canvas = document.createElement('canvas')
  canvas.width = targetSide
  canvas.height = targetSide
  const ctx = canvas.getContext('2d')
  if (!ctx) throw new Error('Canvas недоступен')
  ctx.drawImage(img, sx, sy, side, side, 0, 0, targetSide, targetSide)

  const blob = await new Promise<Blob>((resolve, reject) => {
    canvas.toBlob((b) => (b ? resolve(b) : reject(new Error('Не удалось обработать изображение'))), 'image/jpeg', 0.92)
  })

  return new File([blob], 'avatar.jpg', { type: 'image/jpeg' })
}
