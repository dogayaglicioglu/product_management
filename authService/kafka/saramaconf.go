package kafka

import "github.com/IBM/sarama"

func SaramaConf() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true                //kafkaya mesaj ulaştığında producer bunu bildirir.
	config.Producer.RequiredAcks = sarama.WaitForAll       //gönderilen mesajın tüm replikalar tarafından alınıp onaylanmasını bekler.
	config.Producer.Retry.Max = 3                          //mesajın gönderilmesi sırasında sorun olursa producer 3 kere deneyecek.
	config.Producer.Compression = sarama.CompressionSnappy //performansı arttırmak için mesajların sıkıştırılmasını sağlar.
	return config
}
