package progress

// import (
// 	"time"

// 	"github.com/schollz/progressbar/v3"
// )

// // NewProgressBar 创建一个进度条
// //
// // 参数:
// //   - total: 进度条总大小
// //   - description: 进度条描述信息
// //
// // 返回:
// //   - *progressbar.ProgressBar: 进度条指针
// func NewProgressBar(total int64, description string) *progressbar.ProgressBar {
// 	return progressbar.NewOptions64(
// 		total,                             // 进度条总大小
// 		progressbar.OptionClearOnFinish(), // 完成后清除进度条
// 		progressbar.OptionSetDescription(description), // 进度条描述信息
// 		progressbar.OptionSetElapsedTime(true),        // 显示已用时间
// 		progressbar.OptionSetPredictTime(true),        // 显示预计剩余时间
// 		progressbar.OptionSetRenderBlankState(true),   // 在进度条完成之前显示空白状态
// 		progressbar.OptionShowBytes(true),             // 显示进度条传输的字节
// 		progressbar.OptionShowCount(),                 // 显示当前进度的总和
// 		//progressbar.OptionSetTheme(progressbar.ThemeASCII), // ASCII 进度条主题(默认为 Unicode 进度条主题)
// 	)
// }

// func main() {
// 	// 模拟总任务量为100个单位
// 	total := int64(100)

// 	// 创建进度条，描述信息为"处理中..."
// 	pbar := NewProgressBar(total, "处理中...")

// 	// 模拟任务进度
// 	for i := int64(0); i < total; i++ {
// 		// 每次完成一个单位，更新进度条
// 		pbar.Add(1)

// 		// 模拟处理耗时
// 		time.Sleep(50 * time.Millisecond)
// 	}

// 	// 确保进度条完成
// 	pbar.Finish()

// 	// 关闭进度条
// 	pbar.Close()

// 	// 进度完成后输出提示
// 	println("任务已完成!")
// }
