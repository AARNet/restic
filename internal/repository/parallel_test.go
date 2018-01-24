package repository_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/restic"

	"github.com/restic/restic/internal/repository"
	rtest "github.com/restic/restic/internal/test"
)

type testIDs []string

var lister = testIDs{
	"40bb581cd36de952985c97a3ff6b21df41ee897d4db2040354caa36a17ff5268",
	"2e15811a4d14ffac66d36a9ff456019d8de4c10c949d45b643f8477d17e92ff3",
	"70c11b3ed521ad6b76d905c002ca98b361fca06aca060a063432c7311155a4da",
	"8056a33e75dccdda701b6c989c7ed0cb71bbb6da13c6427fe5986f0896cc91c0",
	"79d8776200596aa0237b10d470f7b850b86f8a1a80988ef5c8bee2874ce992e2",
	"f9f1f29791c6b79b90b35efd083f17a3b163bbbafb1a2fdf43d46d56cffda289",
	"3834178d05d0f6dd07f872ee0262ff1ace0f0f375768227d3c902b0b66591369",
	"66d5cc68c9186414806f366ae5493ce7f229212993750a4992be4030f6af28c5",
	"ebca5af4f397944f68cd215e3dfa2b197a7ba0f7c17d65d9f7390d0a15cde296",
	"d4511ce6ff732d106275a57e40745c599e987c0da44c42cddbef592aac102437",
	"f366202f0bfeefaedd7b49e2f21a90d3cbddb97d257a74d788dd34e19a684dae",
	"a5c17728ab2433cd50636dd5c6c7068c7a44f2999d09c46e8f528466da8a059d",
	"bae0f9492b9b208233029b87692a1a55cbd7fbe1cf3f6d7bc693ac266a6d6f0e",
	"9d500187913c7510d71d1902703d312c7aaa56f1e98351385b9535fdabae595e",
	"ffbddd8a4c1e54d258bb3e16d3929b546b61af63cb560b3e3061a8bef5b24552",
	"201bb3abf655e7ef71e79ed4fb1079b0502b5acb4d9fad5e72a0de690c50a386",
	"08eb57bbd559758ea96e99f9b7688c30e7b3bcf0c4562ff4535e2d8edeffaeed",
	"e50b7223b04985ff38d9e11d1cba333896ef4264f82bd5d0653a028bce70e542",
	"65a9421cd59cc7b7a71dcd9076136621af607fb4701d2e5c2af23b6396cf2f37",
	"995a655b3521c19b4d0c266222266d89c8fc62889597d61f45f336091e646d57",
	"51ec6f0bce77ed97df2dd7ae849338c3a8155a057da927eedd66e3d61be769ad",
	"7b3923a0c0666431efecdbf6cb171295ec1710b6595eebcba3b576b49d13e214",
	"2cedcc3d14698bea7e4b0546f7d5d48951dd90add59e6f2d44b693fd8913717d",
	"fd6770cbd54858fdbd3d7b4239b985e5599180064d93ca873f27e86e8407d011",
	"9edc51d8e6e04d05c9757848c1bfbfdc8e86b6330982294632488922e59fdb1b",
	"1a6c4fbb24ad724c968b2020417c3d057e6c89e49bdfb11d91006def65eab6a0",
	"cb3b29808cd0adfa2dca1f3a04f98114fbccf4eb487cdd4022f49bd70eeb049b",
	"f55edcb40c619e29a20e432f8aaddc83a649be2c2d1941ccdc474cd2af03d490",
	"e8ccc1763a92de23566b95c3ad1414a098016ece69a885fc8a72782a7517d17c",
	"0fe2e3db8c5a12ad7101a63a0fffee901be54319cfe146bead7aec851722f82d",
	"36be45a6ae7c95ad97cee1b33023be324bce7a7b4b7036e24125679dd9ff5b44",
	"1685ed1a57c37859fbef1f7efb7509f20b84ec17a765605de43104d2fa37884b",
	"9d83629a6a004c505b100a0b5d0b246833b63aa067aa9b59e3abd6b74bc4d3a8",
	"be49a66b60175c5e2ee273b42165f86ef11bb6518c1c79950bcd3f4c196c98bd",
	"0fd89885d821761b4a890782908e75793028747d15ace3c6cbf0ad56582b4fa5",
	"94a767519a4e352a88796604943841fea21429f3358b4d5d55596dbda7d15dce",
	"8dd07994afe6e572ddc9698fb0d13a0d4c26a38b7992818a71a99d1e0ac2b034",
	"f7380a6f795ed31fbeb2945c72c5fd1d45044e5ab152311e75e007fa530f5847",
	"5ca1ce01458e484393d7e9c8af42b0ff37a73a2fee0f18e14cff0fb180e33014",
	"8f44178be3fe0a2bd41f922576fb7a9b19d589754504be746f56c759df328fda",
	"12d33847c2be711c989f37360dd7aa8537fd14972262a4530634a08fdf32a767",
	"31e077f5080f78846a00093caff2b6b839519cc47516142eeba9c41d4072a605",
	"14f01db8a0054e70222b76d2555d70114b4bf8a0f02084324af2df226f14a795",
	"7f5dbbaf31b4551828e8e76cef408375db9fbcdcdb6b5949f2d1b0c4b8632132",
	"42a5d9b9bb7e4a16f23ba916bcf87f38c1aa1f2de2ab79736f725850a8ff6a1b",
	"e06f8f901ea708beba8712a11b6e2d0be7c4b018d0254204ef269bcdf5e8c6cc",
	"d9ba75785bf45b0c4fd3b2365c968099242483f2f0d0c7c20306dac11fae96e9",
	"428debbb280873907cef2ec099efe1566e42a59775d6ec74ded0c4048d5a6515",
	"3b51049d4dae701098e55a69536fa31ad2be1adc17b631a695a40e8a294fe9c0",
	"168f88aa4b105e9811f5f79439cc1a689be4eec77f3361d42f22fe8f7ddc74a9",
	"0baa0ab2249b33d64449a899cb7bd8eae5231f0d4ff70f09830dc1faa2e4abee",
	"0c3896d346b580306a49de29f3a78913a41e14b8461b124628c33a64636241f2",
	"b18313f1651c15e100e7179aa3eb8ffa62c3581159eaf7f83156468d19781e42",
	"996361f7d988e48267ccc7e930fed4637be35fe7562b8601dceb7a32313a14c8",
	"dfb4e6268437d53048d22b811048cd045df15693fc6789affd002a0fc80a6e60",
	"34dd044c228727f2226a0c9c06a3e5ceb5e30e31cb7854f8fa1cde846b395a58",
}

func (tests testIDs) List(ctx context.Context, t restic.FileType, fn func(restic.FileInfo) error) error {
	for i := 0; i < 500; i++ {
		for _, id := range tests {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			fi := restic.FileInfo{
				Name: id,
			}

			err := fn(fi)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestFilesInParallel(t *testing.T) {
	f := func(ctx context.Context, id string) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	for n := 1; n < 5; n++ {
		err := repository.FilesInParallel(context.TODO(), lister, restic.DataFile, n*100, f)
		rtest.OK(t, err)
	}
}

var errTest = errors.New("test error")

func TestFilesInParallelWithError(t *testing.T) {
	f := func(ctx context.Context, id string) error {
		time.Sleep(1 * time.Millisecond)

		if rand.Float32() < 0.01 {
			return errTest
		}

		return nil
	}

	for n := 1; n < 5; n++ {
		err := repository.FilesInParallel(context.TODO(), lister, restic.DataFile, n*100, f)
		if err != errTest {
			t.Fatalf("wrong error returned, want %q, got %v", errTest, err)
		}
	}
}
