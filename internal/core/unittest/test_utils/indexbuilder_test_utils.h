// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License

#pragma once

#include <tuple>
#include <map>
#include <limits>
#include <cmath>
#include <google/protobuf/text_format.h>

#include "DataGen.h"
#include "index/ScalarIndex.h"
#include "index/StringIndex.h"
#include "indexbuilder/ScalarIndexCreator.h"
#include "indexbuilder/VecIndexCreator.h"
#include "indexbuilder/helper.h"
#include "indexbuilder/index_c.h"
#include "indexbuilder/utils.h"
#include "knowhere/index/VecIndexFactory.h"
#include "knowhere/index/vector_index/helpers/IndexParameter.h"
#include "knowhere/index/vector_index/adapter/VectorAdapter.h"
#include "pb/index_cgo_msg.pb.h"

constexpr int64_t DIM = 8;
constexpr int64_t NQ = 10;
constexpr int64_t K = 4;
#ifdef MILVUS_GPU_VERSION
int DEVICEID = 0;
#endif

namespace indexcgo = milvus::proto::indexcgo;
namespace schemapb = milvus::proto::schema;
using milvus::indexbuilder::MapParams;
using milvus::indexbuilder::ScalarIndexCreator;
using ScalarTestParams = std::pair<MapParams, MapParams>;
using milvus::scalar::ScalarIndexPtr;
using milvus::scalar::StringIndexPtr;

namespace {
auto
generate_conf(const knowhere::IndexType& index_type, const knowhere::MetricType& metric_type) {
    if (index_type == knowhere::IndexEnum::INDEX_FAISS_IDMAP) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_FAISS_IVFPQ) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::NLIST, 16},
            {knowhere::indexparam::NPROBE, 4},
            {knowhere::indexparam::M, 4},
            {knowhere::indexparam::NBITS, 8},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_FAISS_IVFFLAT) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::NLIST, 16},
            {knowhere::indexparam::NPROBE, 4},
#ifdef MILVUS_GPU_VERSION
            {knowhere::meta::DEVICE_ID, DEVICEID},
#endif
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_FAISS_IVFSQ8) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::NLIST, 16},
            {knowhere::indexparam::NPROBE, 4},
            {knowhere::indexparam::NBITS, 8},
#ifdef MILVUS_GPU_VERSION
            {knowhere::meta::DEVICE_ID, DEVICEID},
#endif
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_FAISS_BIN_IVFFLAT) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::NLIST, 16},
            {knowhere::indexparam::NPROBE, 4},
            {knowhere::indexparam::M, 4},
            {knowhere::indexparam::NBITS, 8},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_FAISS_BIN_IDMAP) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
        };
#ifdef MILVUS_SUPPORT_NSG
    } else if (index_type == knowhere::IndexEnum::INDEX_NSG) {
        return knowhere::Config{
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::NLIST, 163},
            {knowhere::indexparam::NPROBE, 8},
            {knowhere::indexparam::KNNG, 20},
            {knowhere::indexparam::SEARCH_LENGTH, 40},
            {knowhere::indexparam::OUT_DEGREE, 30},
            {knowhere::indexparam::CANDIDATE, 100},
        };
#endif
#ifdef MILVUS_SUPPORT_SPTAG
    } else if (index_type == knowhere::IndexEnum::INDEX_SPTAG_KDT_RNT) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, 10},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_SPTAG_BKT_RNT) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, 10},
        };
#endif
    } else if (index_type == knowhere::IndexEnum::INDEX_HNSW) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::HNSW_M, 16},
            {knowhere::indexparam::EFCONSTRUCTION, 200},
            {knowhere::indexparam::EF, 200},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_ANNOY) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::N_TREES, 4},
            {knowhere::indexparam::SEARCH_K, 100},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_RHNSWFlat) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::HNSW_M, 16},
            {knowhere::indexparam::EFCONSTRUCTION, 200},
            {knowhere::indexparam::EF, 200},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_RHNSWPQ) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::HNSW_M, 16},
            {knowhere::indexparam::EFCONSTRUCTION, 200},
            {knowhere::indexparam::EF, 200},
            {knowhere::indexparam::PQ_M, 8},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_RHNSWSQ) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::HNSW_M, 16},
            {knowhere::indexparam::EFCONSTRUCTION, 200},
            {knowhere::indexparam::EF, 200},
        };
#ifdef MILVUS_SUPPORT_NGT
    } else if (index_type == knowhere::IndexEnum::INDEX_NGTPANNG) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::EDGE_SIZE, 10},
            {knowhere::indexparam::EPSILON, 0.1},
            {knowhere::indexparam::MAX_SEARCH_EDGES, 50},
            {knowhere::indexparam::FORCEDLY_PRUNED_EDGE_SIZE, 60},
            {knowhere::indexparam::SELECTIVELY_PRUNED_EDGE_SIZE, 30},
        };
    } else if (index_type == knowhere::IndexEnum::INDEX_NGTONNG) {
        return knowhere::Config{
            {knowhere::meta::SLICE_SIZE, knowhere::index_file_slice_size},
            {knowhere::meta::METRIC_TYPE, metric_type},
            {knowhere::meta::DIM, DIM},
            {knowhere::meta::TOPK, K},
            {knowhere::indexparam::EDGE_SIZE, 20},
            {knowhere::indexparam::EPSILON, 0.1},
            {knowhere::indexparam::MAX_SEARCH_EDGES, 50},
            {knowhere::indexparam::OUTGOING_EDGE_SIZE, 5},
            {knowhere::indexparam::INCOMING_EDGE_SIZE, 40},
        };
#endif
    }
    return knowhere::Config();
}

auto
generate_params(const knowhere::IndexType& index_type, const knowhere::MetricType& metric_type) {
    namespace indexcgo = milvus::proto::indexcgo;

    indexcgo::TypeParams type_params;
    indexcgo::IndexParams index_params;

    auto configs = generate_conf(index_type, metric_type);
    for (auto& [key, value] : configs.items()) {
        auto param = index_params.add_params();
        auto value_str = value.is_string() ? value.get<std::string>() : value.dump();
        param->set_key(key);
        param->set_value(value_str);
    }

    auto param = index_params.add_params();
    param->set_key("index_type");
    param->set_value(std::string(index_type));

    return std::make_tuple(type_params, index_params);
}

auto
GenDataset(int64_t N, const knowhere::MetricType& metric_type, bool is_binary, int64_t dim = DIM) {
    auto schema = std::make_shared<milvus::Schema>();
    if (!is_binary) {
        schema->AddDebugField("fakevec", milvus::DataType::VECTOR_FLOAT, dim, metric_type);
        return milvus::segcore::DataGen(schema, N);
    } else {
        schema->AddDebugField("fakebinvec", milvus::DataType::VECTOR_BINARY, dim, metric_type);
        return milvus::segcore::DataGen(schema, N);
    }
}

using QueryResultPtr = std::unique_ptr<milvus::indexbuilder::VecIndexCreator::QueryResult>;
void
PrintQueryResult(const QueryResultPtr& result) {
    auto nq = result->nq;
    auto k = result->topk;

    std::stringstream ss_id;
    std::stringstream ss_dist;

    for (auto i = 0; i < nq; i++) {
        for (auto j = 0; j < k; ++j) {
            ss_id << result->ids[i * k + j] << " ";
            ss_dist << result->distances[i * k + j] << " ";
        }
        ss_id << std::endl;
        ss_dist << std::endl;
    }
    std::cout << "id\n" << ss_id.str() << std::endl;
    std::cout << "dist\n" << ss_dist.str() << std::endl;
}

float
L2(const float* point_a, const float* point_b, int dim) {
    float dis = 0;
    for (auto i = 0; i < dim; i++) {
        auto c_a = point_a[i];
        auto c_b = point_b[i];
        dis += pow(c_b - c_a, 2);
    }
    return dis;
}

int
hamming_weight(uint8_t n) {
    int count = 0;
    while (n != 0) {
        count += n & 1;
        n >>= 1;
    }
    return count;
}
float
Jaccard(const uint8_t* point_a, const uint8_t* point_b, int dim) {
    float dis;
    int len = dim / 8;
    float intersection = 0;
    float union_num = 0;
    for (int i = 0; i < len; i++) {
        intersection += hamming_weight(point_a[i] & point_b[i]);
        union_num += hamming_weight(point_a[i] | point_b[i]);
    }
    dis = 1 - (intersection / union_num);
    return dis;
}

float
CountDistance(
    const void* point_a, const void* point_b, int dim, const knowhere::MetricType& metric, bool is_binary = false) {
    if (point_a == nullptr || point_b == nullptr) {
        return std::numeric_limits<float>::max();
    }
    if (metric == knowhere::metric::L2) {
        return L2(static_cast<const float*>(point_a), static_cast<const float*>(point_b), dim);
    } else if (metric == knowhere::metric::JACCARD) {
        return Jaccard(static_cast<const uint8_t*>(point_a), static_cast<const uint8_t*>(point_b), dim);
    } else {
        return std::numeric_limits<float>::max();
    }
}

void
CheckDistances(const QueryResultPtr& result,
               const knowhere::DatasetPtr& base_dataset,
               const knowhere::DatasetPtr& query_dataset,
               const knowhere::MetricType& metric,
               const float threshold = 1.0e-5) {
    auto base_vecs = (float*)knowhere::GetDatasetTensor(base_dataset);
    auto query_vecs = (float*)knowhere::GetDatasetTensor(query_dataset);
    auto dim = knowhere::GetDatasetDim(base_dataset);
    auto nq = result->nq;
    auto k = result->topk;
    for (auto i = 0; i < nq; i++) {
        for (auto j = 0; j < k; ++j) {
            auto dis = result->distances[i * k + j];
            auto id = result->ids[i * k + j];
            auto count_dis = CountDistance(query_vecs + i * dim, base_vecs + id * dim, dim, metric);
            // assert(std::abs(dis - count_dis) < threshold);
        }
    }
}

auto
generate_type_params(const MapParams& m) {
    indexcgo::TypeParams p;
    for (const auto& [k, v] : m) {
        auto kv = p.add_params();
        kv->set_key(k);
        kv->set_value(v);
    }
    std::string str;
    auto ok = google::protobuf::TextFormat::PrintToString(p, &str);
    Assert(ok);
    return str;
}

auto
generate_index_params(const MapParams& m) {
    indexcgo::IndexParams p;
    for (const auto& [k, v] : m) {
        auto kv = p.add_params();
        kv->set_key(k);
        kv->set_value(v);
    }
    std::string str;
    auto ok = google::protobuf::TextFormat::PrintToString(p, &str);
    Assert(ok);
    return str;
}

// TODO: std::is_arithmetic_v, hard to compare float point value. std::is_integral_v.
template <typename T, typename = typename std::enable_if_t<std::is_arithmetic_v<T> || std::is_same_v<T, std::string>>>
inline std::vector<T>
GenArr(int64_t n) {
    auto max_i8 = std::numeric_limits<int8_t>::max() - 1;
    std::vector<T> arr;
    arr.resize(n);
    for (int64_t i = 0; i < n; i++) {
        arr[i] = static_cast<T>(rand() % max_i8);
    }
    std::sort(arr.begin(), arr.end());
    return arr;
}

inline auto
GenStrArr(int64_t n) {
    using T = std::string;
    std::vector<T> arr;
    arr.resize(n);
    for (int64_t i = 0; i < n; i++) {
        auto gen = std::to_string(std::rand());
        arr[i] = gen;
    }
    std::sort(arr.begin(), arr.end());
    return arr;
}

template <>
inline std::vector<std::string>
GenArr<std::string>(int64_t n) {
    return GenStrArr(n);
}

std::vector<ScalarTestParams>
GenBoolParams() {
    std::vector<ScalarTestParams> ret;
    ret.emplace_back(ScalarTestParams(MapParams(), {{"index_type", "inverted_index"}}));
    ret.emplace_back(ScalarTestParams(MapParams(), {{"index_type", "flat"}}));
    return ret;
}

std::vector<ScalarTestParams>
GenStringParams() {
    std::vector<ScalarTestParams> ret;
    ret.emplace_back(ScalarTestParams(MapParams(), {{"index_type", "marisa"}}));
    return ret;
}

template <typename T, typename = typename std::enable_if_t<std::is_arithmetic_v<T> | std::is_same_v<std::string, T>>>
inline std::vector<ScalarTestParams>
GenParams() {
    if (std::is_same_v<std::string, T>) {
        return GenStringParams();
    }

    if (std::is_same_v<T, bool>) {
        return GenBoolParams();
    }

    std::vector<ScalarTestParams> ret;
    ret.emplace_back(ScalarTestParams(MapParams(), {{"index_type", "inverted_index"}}));
    ret.emplace_back(ScalarTestParams(MapParams(), {{"index_type", "flat"}}));
    return ret;
}

void
PrintMapParam(const ScalarTestParams& tp) {
    for (const auto& [k, v] : tp.first) {
        std::cout << "k: " << k << ", v: " << v << std::endl;
    }
    for (const auto& [k, v] : tp.second) {
        std::cout << "k: " << k << ", v: " << v << std::endl;
    }
}

void
PrintMapParams(const std::vector<ScalarTestParams>& tps) {
    for (const auto& tp : tps) {
        PrintMapParam(tp);
    }
}

// memory generated by this function should be freed by the caller.
auto
GenDsFromPB(const google::protobuf::Message& msg) {
    auto data = new char[msg.ByteSize()];
    msg.SerializeToArray(data, msg.ByteSize());
    return knowhere::GenDataset(msg.ByteSize(), 8, data);
}

template <typename T>
inline std::vector<std::string>
GetIndexTypes() {
    return std::vector<std::string>{"inverted_index"};
}

template <>
inline std::vector<std::string>
GetIndexTypes<std::string>() {
    return std::vector<std::string>{"marisa"};
}

}  // namespace
